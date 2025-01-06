package app

import (
	"github.com/spf13/viper"

	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
)

type HttpConfig struct {
	Host string
	Port uint16
}

type Config struct {
	App  AppConfig
	DB   repository.Config
	Http HttpConfig
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
	Name string
}
