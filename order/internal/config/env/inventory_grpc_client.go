package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type inventoryGrpcClientEnvConfig struct {
	GrpcHost string `env:"INVENTORY_GRPC_HOST,required"`
	GrpcPort int    `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryGrpcClientConfig struct {
	raw inventoryGrpcClientEnvConfig
}

func NewInventoryGrpcClientConfig() (*inventoryGrpcClientConfig, error) {
	var raw inventoryGrpcClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &inventoryGrpcClientConfig{raw: raw}, nil
}

func (invGrpc *inventoryGrpcClientConfig) Host() string {
	return invGrpc.raw.GrpcHost
}

func (invGrpc *inventoryGrpcClientConfig) Port() int {
	return invGrpc.raw.GrpcPort
}

func (invGrpc *inventoryGrpcClientConfig) Address() string {
	return fmt.Sprintf("%s:%d", invGrpc.raw.GrpcHost, invGrpc.raw.GrpcPort)
}
