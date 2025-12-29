package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type paymentGrpcClientEnvConfig struct {
	GrpcHost string `env:"PAYMENT_GRPC_HOST,required"`
	GrpcPort int    `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentGrpcClientConfig struct {
	raw paymentGrpcClientEnvConfig
}

func NewPaymentGrpcClientConfig() (*paymentGrpcClientConfig, error) {
	var raw paymentGrpcClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &paymentGrpcClientConfig{raw: raw}, nil
}

func (pc *paymentGrpcClientConfig) Host() string {
	return pc.raw.GrpcHost
}
func (pc *paymentGrpcClientConfig) Port() int {
	return pc.raw.GrpcPort
}
func (pc *paymentGrpcClientConfig) Address() string {
	return fmt.Sprintf("%s:%d", pc.raw.GrpcHost, pc.raw.GrpcPort)
}