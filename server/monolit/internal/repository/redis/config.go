package redis

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DB       int
}

func (c Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
