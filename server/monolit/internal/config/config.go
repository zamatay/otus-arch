package config

import (
	"github.com/spf13/viper"

	"githib.com/zamatay/otus/arch/lesson-1/internal/api"
	"githib.com/zamatay/otus/arch/lesson-1/internal/grpconnection"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

type Config struct {
	App   AppConfig
	DB    map[string][]*repository.Config
	Cache redis.Config
	Http  api.Config
	Kafka kafka.Config
	GRPC  grpcserver.Config
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

type AppConfig struct {
	Name   string
	Secret string
}
