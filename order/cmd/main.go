package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/delyke/go_workspace_example/order/internal/api/order/v1"
	inventoryClient "github.com/delyke/go_workspace_example/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/delyke/go_workspace_example/order/internal/client/grpc/payment/v1"
	"github.com/delyke/go_workspace_example/order/internal/migrator"
	orderRepo "github.com/delyke/go_workspace_example/order/internal/repository/order"
	orderService "github.com/delyke/go_workspace_example/order/internal/service/order"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	inventoryAddress  = "localhost:50051"
	paymentAddress    = "localhost:50052"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
		return
	}

	dbURI := os.Getenv("DB_URI")

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("failed to ping database: %v", err)
		return
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	poolCfg, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*poolCfg.ConnConfig), migrationsDir)
	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	connInventory, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("inventory grpc connect failed: %v", err)
	}
	defer func() {
		if cerr := connInventory.Close(); cerr != nil {
			log.Printf("inventory grpc close failed: %v", cerr)
		}
	}()

	generatedInvClient := inventoryV1.NewInventoryServiceClient(connInventory)

	invClient := inventoryClient.NewClient(generatedInvClient)

	connPayment, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to create payment client: %v", err)
		return
	}
	defer func() {
		if cerr := connPayment.Close(); cerr != nil {
			log.Printf("Failed to close payment client: %v", cerr)
		}
	}()

	generatedPaymentClient := paymentV1.NewPaymentServiceClient(connPayment)
	payClient := paymentClient.NewClient(generatedPaymentClient)

	repo := orderRepo.NewRepository(pool)
	service := orderService.NewService(repo, invClient, payClient)
	api := orderV1Api.NewApi(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
