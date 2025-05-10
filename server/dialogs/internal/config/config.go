package config

import (
	"github.com/spf13/viper"

	"dialogs/internal/api"
	"dialogs/internal/api/grpcclient"
	"dialogs/internal/repository"
)

type AppConfig struct {
	Name   string
	Secret string
}

type Config struct {
	App AppConfig
	DB  *repository.Config
	//Cache redis.Config
	Http api.Config
	GRPC grpcclient.Config
	//Kafka kafka.Config
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
