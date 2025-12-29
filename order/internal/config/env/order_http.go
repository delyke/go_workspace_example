package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
)

type httpEnvConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
	ReadTimeout string `env:"HTTP_READ_TIMEOUT,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &httpConfig{raw: raw}, nil
}

func (h *httpConfig) Host() string {
	return h.raw.Host
}
func (h *httpConfig) Port() string {
	return h.raw.Port
}

func (h *httpConfig) Address() string {
	return fmt.Sprintf("%s:%s", h.raw.Host, h.raw.Port)
}

func (h *httpConfig) ReadTimeout() time.Duration {
	readTimeout, err := time.ParseDuration(h.raw.ReadTimeout)
	if err != nil {
		return 5 * time.Second
	}
	return readTimeout
}