package config

import (
	"github.com/joho/godotenv"

	"github.com/delyke/go_workspace_example/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Mongo         MongoConfig
	InventoryGRPC InventoryConfig
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

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := env.NewInventoryGrpcConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Logger:        loggerCfg,
		Mongo:         mongoCfg,
		InventoryGRPC: inventoryCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
