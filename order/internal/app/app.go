package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/delyke/go_workspace_example/order/internal/config"
	"github.com/delyke/go_workspace_example/platform/pkg/closer"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
	"github.com/delyke/go_workspace_example/platform/pkg/migrator"
)

type App struct {
	diContainer *diContainer
	router      *chi.Mux
	httpServer  *http.Server
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
	return a.runHttpServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx2 context.Context) error{
		a.initDI,
		a.initLogger,
		a.initMigrator,
		a.initCloser,
		a.initRouter,
		a.initServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initMigrator(ctx context.Context) error {
	poolCfg, err := pgxpool.ParseConfig(config.AppConfig().Postgres.URI())
	if err != nil {
		return err
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*poolCfg.ConnConfig), config.AppConfig().Postgres.MigrationDirectory())
	err = migratorRunner.Up(ctx)
	if err != nil {
		logger.Error(ctx, "migrator up failed", zap.Error(err))
		return err
	}
	logger.Info(ctx, "migrator up done")
	return nil
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", a.diContainer.OrderServer(ctx))
	a.router = r
	return nil
}

func (a *App) initServer(_ context.Context) error {
	server := &http.Server{
		Addr:              net.JoinHostPort(config.AppConfig().HTTP.Host(), config.AppConfig().HTTP.Port()),
		Handler:           a.router,
		ReadHeaderTimeout: config.AppConfig().HTTP.ReadTimeout(), // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}
	a.httpServer = server
	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, "starting http server", zap.String("host", config.AppConfig().HTTP.Host()), zap.String("port", config.AppConfig().HTTP.Port()))
	closer.AddNamed("Http Server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
