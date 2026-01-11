package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host               string `env:"POSTGRES_HOST,required"`
	Port               int    `env:"EXTERNAL_POSTGRES_PORT,required"`
	User               string `env:"POSTGRES_USER,required"`
	Password           string `env:"POSTGRES_PASSWORD,required"`
	Database           string `env:"POSTGRES_DB,required"`
	MigrationDirectory string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &postgresConfig{raw: raw}, nil
}

func (c *postgresConfig) Host() string {
	return c.raw.Host
}

func (c *postgresConfig) Port() int {
	return c.raw.Port
}

func (c *postgresConfig) User() string {
	return c.raw.User
}

func (c *postgresConfig) Password() string {
	return c.raw.Password
}

func (c *postgresConfig) Database() string {
	return c.raw.Database
}

func (c *postgresConfig) MigrationDirectory() string {
	return c.raw.MigrationDirectory
}

func (c *postgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.raw.User,
		c.raw.Password,
		c.raw.Host,
		c.raw.Port,
		c.raw.Database,
	)
}
