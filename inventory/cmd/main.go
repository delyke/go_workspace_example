package main

import (
	"context"
	"fmt"
	"github.com/delyke/go_workspace_example/inventory/internal/config"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/delyke/go_workspace_example/inventory/internal/api/inventory/v1"
	partRepository "github.com/delyke/go_workspace_example/inventory/internal/repository/part"
	partService "github.com/delyke/go_workspace_example/inventory/internal/service/part"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", config.AppConfig().InventoryGRPC.Address()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("Error connecting to mongo db: %v", err)
	}

	defer func() {
		cerr := mongoClient.Disconnect(ctx)
		if cerr != nil {
			log.Printf("Error disconnecting from mongo db: %v", cerr)
		}
	}()

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging mongo db: %v", err)
		return
	}

	db := mongoClient.Database("inventory")

	partRepo := partRepository.NewRepository(db)
	partSrv := partService.NewService(partRepo)
	invApi := inventoryApi.NewApi(partSrv)

	inventoryV1.RegisterInventoryServiceServer(s, invApi)
	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on %s\n", config.AppConfig().InventoryGRPC.Address())
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down gRPC server...")
	s.GracefulStop()
	log.Printf("Server gracefully stopped")
}
