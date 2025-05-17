package repository

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

func (c Config) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}
