package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCConfig interface {
	Host() string
	Port() string
	Address() string
}
