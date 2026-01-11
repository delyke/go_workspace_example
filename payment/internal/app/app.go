package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/delyke/go_workspace_example/payment/internal/config"
	"github.com/delyke/go_workspace_example/platform/pkg/closer"
	"github.com/delyke/go_workspace_example/platform/pkg/grpc/health"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		logger.Error(ctx, "failed to listen", zap.Error(err))
		return err
	}

	closer.AddNamed("TCP Listener", func(ctx context.Context) error {
		lerr := lis.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			logger.Error(ctx, "failed to close listener", zap.Error(lerr))
			return lerr
		}
		return nil
	})

	a.listener = lis
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("GRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})
	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentAPIV1(ctx))
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on  %s", config.AppConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		logger.Error(ctx, "failed to serve grpc", zap.Error(err))
		return err
	}
	return nil
}
