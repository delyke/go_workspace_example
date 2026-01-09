package config

import (
	"github.com/delyke/go_workspace_example/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	paymentGRPCConf, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: paymentGRPCConf,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
