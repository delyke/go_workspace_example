package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/delyke/go_workspace_example/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger          LoggerConfig
	Postgres        PostgresConfig
	HTTP            OrderHTTPConfig
	PaymentClient   PaymentGRPCClientConfig
	InventoryClient InventoryGRPCClientConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}
	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}
	paymentCfg, err := env.NewPaymentGrpcClientConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := env.NewInventoryGrpcClientConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:          loggerCfg,
		Postgres:        postgresCfg,
		HTTP:            httpCfg,
		PaymentClient:   paymentCfg,
		InventoryClient: inventoryCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
