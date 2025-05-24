package config

import (
	"github.com/spf13/viper"

	"counter/internal/api"
	"counter/internal/api/grpcclient"
	"counter/internal/repository/redis"
)

type AppConfig struct {
	Name   string
	Secret string
}

type Config struct {
	App        AppConfig
	Cache      redis.Config
	Http       api.Config
	GRPCClient grpcclient.Config
	GRPCServer grpcclient.ConfigServer
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
