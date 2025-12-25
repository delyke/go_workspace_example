package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApiV1 "github.com/delyke/go_workspace_example/payment/internal/api/payment/v1"
	paymentService "github.com/delyke/go_workspace_example/payment/internal/service/payment"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	pService := paymentService.NewService()
	apiV1 := paymentApiV1.NewApi(pService)

	paymentV1.RegisterPaymentServiceServer(s, apiV1)
	reflection.Register(s)

	go func() {
		log.Printf("payment gRPC server listening at %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
