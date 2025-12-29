package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type inventoryGrpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGrpcConfig struct {
	raw inventoryGrpcEnvConfig
}

func NewInventoryGrpcConfig() (*inventoryGrpcConfig, error) {
	var raw inventoryGrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &inventoryGrpcConfig{raw: raw}, nil
}

func (i *inventoryGrpcConfig) Host() string {
	return i.raw.Host
}
func (i *inventoryGrpcConfig) Port() string {
	return i.raw.Port
}

func (i *inventoryGrpcConfig) Address() string {
	return fmt.Sprintf("%s:%s", i.Host(), i.Port())
}