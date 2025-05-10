package grpcserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/zamatay/otus/arch/lesson-1/pkg/api"
)

type AddHandler interface {
	AddHandle(path string, handler http.Handler)
}

type Service struct {
	pb.UnimplementedMonolitServer
	listener net.Listener
	server   *grpc.Server
	router   *net.Listener
}

func NewGRPCServer(cfg *Config) (*Service, error) {
	srv := new(Service)
	var err error
	srv.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	srv.server = grpc.NewServer()

	reflection.Register(srv.server)

	return srv, nil
}

func (s *Service) Start() error {
	return s.server.Serve(s.listener)
}

func (s *Service) Register(ctx context.Context, service AddHandler) error {
	pb.RegisterMonolitServer(s.server, s)
	mux := runtime.NewServeMux()
	if err := pb.RegisterMonolitHandlerServer(ctx, mux, s); err != nil {
		return err
	}
	service.AddHandle("/", mux)
	return nil
}

func (s *Service) HealthCheck(ctx context.Context, empty *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "ok"}, nil
}
