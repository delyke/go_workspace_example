package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1Api "github.com/delyke/go_workspace_example/order/internal/api/order/v1"
	grpcClients "github.com/delyke/go_workspace_example/order/internal/client/grpc"
	inventoryClient "github.com/delyke/go_workspace_example/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/delyke/go_workspace_example/order/internal/client/grpc/payment/v1"
	"github.com/delyke/go_workspace_example/order/internal/config"
	"github.com/delyke/go_workspace_example/order/internal/repository"
	orderRepo "github.com/delyke/go_workspace_example/order/internal/repository/order"
	"github.com/delyke/go_workspace_example/order/internal/service"
	orderService "github.com/delyke/go_workspace_example/order/internal/service/order"
	"github.com/delyke/go_workspace_example/platform/pkg/closer"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API          orderV1.Handler
	orderService        service.OrderService
	orderRepository     repository.OrderRepository
	paymentGRPCClient   grpcClients.PaymentClient
	inventoryGRPCClient grpcClients.InventoryClient
	pool                *pgxpool.Pool
	poolCfg             *pgxpool.Config
	orderServer         *orderV1.Server
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) OrderServer(ctx context.Context) *orderV1.Server {
	if d.orderServer == nil {
		orderServer, err := orderV1.NewServer(d.OrderV1API(ctx))
		if err != nil {
			panic(fmt.Errorf("failed to initialize order server: %w", err))
		}
		d.orderServer = orderServer
	}
	return d.orderServer
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1Api.NewApi(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepo.NewRepository(d.PostgresPool(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.InventoryGRPCClient(ctx), d.PaymentGRPCClient(ctx))
	}
	return d.orderService
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.pool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}
		closer.AddNamed("Postgres Pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})
		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping postgres: %w", err))
		}
		d.pool = pool
	}
	return d.pool
}

func (d *diContainer) PoolCfg(_ context.Context) *pgxpool.Config {
	if d.poolCfg == nil {
		poolCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}
		d.poolCfg = poolCfg
	}
	return d.poolCfg
}

func (d *diContainer) InventoryGRPCClient(ctx context.Context) grpcClients.InventoryClient {
	if d.inventoryGRPCClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "inventory grpc connect failed", zap.Error(err))
			return nil
		}
		closer.AddNamed("Inventory GRPC Client", func(ctx context.Context) error {
			return conn.Close()
		})
		generatedClient := inventoryV1.NewInventoryServiceClient(conn)
		d.inventoryGRPCClient = inventoryClient.NewClient(generatedClient)
	}
	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(ctx context.Context) grpcClients.PaymentClient {
	if d.paymentGRPCClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "Failed to create payment client:", zap.Error(err))
			return nil
		}
		closer.AddNamed("Payment GRPC Client", func(ctx context.Context) error {
			return conn.Close()
		})

		generatedClient := paymentV1.NewPaymentServiceClient(conn)
		d.paymentGRPCClient = paymentClient.NewClient(generatedClient)
	}
	return d.paymentGRPCClient
}
