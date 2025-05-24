package config

import (
	"github.com/spf13/viper"

	"github.com/zamatay/otus/arch/lesson-1/internal/api"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/grpcclient"
	"github.com/zamatay/otus/arch/lesson-1/internal/grpcserver"
	"github.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository"
	"github.com/zamatay/otus/arch/lesson-1/internal/repository/redis"
)

type Config struct {
	App         AppConfig
	DB          map[string][]*repository.Config
	Cache       redis.Config
	Http        api.Config
	Kafka       kafka.Config
	GRPC        grpcserver.Config
	GRPCCounter grpcclient.Config
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
