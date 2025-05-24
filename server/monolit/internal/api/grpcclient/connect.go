package grpcclient

import (
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateConnection(cfg Config) (*grpc.ClientConn, error) {
	options := []grpc.DialOption{}
	if cfg.UseTls {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), options...)
}
