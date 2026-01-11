package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryConfig interface {
	Host() string
	Port() string
	Address() string
}

type MongoConfig interface {
	ImageName() string
	ExternalPort() string
	Host() string
	Port() string
	Database() string
	AuthDB() string
	User() string
	Password() string
	URI() string
}
