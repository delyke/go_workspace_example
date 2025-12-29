package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	MigrationDirectory() string
	URI() string
}

type OrderHTTPConfig interface {
	Host() string
	Port() string
	Address() string
	ReadTimeout() time.Duration
}

type PaymentGRPCClientConfig interface {
	Host() string
	Port() int
	Address() string
}

type InventoryGRPCClientConfig interface {
	Host() string
	Port() int
	Address() string
}
