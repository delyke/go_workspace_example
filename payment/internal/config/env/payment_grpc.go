package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type paymentGrpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentGrpcConfig struct {
	raw paymentGrpcEnvConfig
}

func NewPaymentGrpcConfig() (*paymentGrpcConfig, error) {
	var raw paymentGrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &paymentGrpcConfig{raw: raw}, nil
}

func (pc *paymentGrpcConfig) Host() string {
	return pc.raw.Host
}
func (pc *paymentGrpcConfig) Port() string {
	return pc.raw.Port
}
func (pc *paymentGrpcConfig) Address() string {
	return fmt.Sprintf("%s:%s", pc.raw.Host, pc.raw.Port)
}