package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/delyke/go_workspace_example/inventory/internal/api/inventory/v1"
	partRepository "github.com/delyke/go_workspace_example/inventory/internal/repository/part"
	partService "github.com/delyke/go_workspace_example/inventory/internal/service/part"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	dbURI := os.Getenv("MONGO_URI")

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
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
		log.Printf("starting gRPC server on port %d\n", grpcPort)
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
