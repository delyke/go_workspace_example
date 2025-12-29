package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type mongoConfigEnv struct {
	ImageName string `env:"MONGO_IMAGE_NAME,required"`
	ExternalPort string `env:"EXTERNAL_MONGO_PORT,required"`
	Host string `env:"MONGO_HOST,required"`
	Port string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	AuthDB string `env:"MONGO_AUTH_DB,required"`
	User string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	raw mongoConfigEnv
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &mongoConfig{raw: raw}, nil
}

func (mc *mongoConfig) ImageName() string  {
	return mc.raw.ImageName
}
func (mc *mongoConfig) ExternalPort() string {
	return mc.raw.ExternalPort
}
func (mc *mongoConfig) Host() string {
	return mc.raw.Host
}
func (mc *mongoConfig) Port() string {
	return mc.raw.Port
}
func (mc *mongoConfig) Database() string {
	return mc.raw.Database
}
func (mc *mongoConfig) AuthDB() string {
	return mc.raw.AuthDB
}
func (mc *mongoConfig) User() string {
	return mc.raw.User
}
func (mc *mongoConfig) Password() string {
	return mc.raw.Password
}
func (mc *mongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		mc.User(),
		mc.Password(),
		mc.Host(),
		mc.Port(),
		mc.Database(),
	)
}
