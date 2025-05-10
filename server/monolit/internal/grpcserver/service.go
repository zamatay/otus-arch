package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"githib.com/zamatay/otus/arch/lesson-1"
)

type Service struct {
	listener net.Listener
	server   *grpc.Server
	router   *net.Listener
}

func NewGRPCServer(cfg Config) (*Service, error) {
	srv := new(Service)
	var err error
	srv.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	srv.server = grpc.NewServer()

	return srv, nil
}

func (s *Service) Start() error {
	return s.server.Serve(s.listener)
}

func (s *Service) Register() error {
	RegisterMonolitServer(s.server, &pb.MonolitServer{})
}
