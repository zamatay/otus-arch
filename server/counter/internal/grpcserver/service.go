package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"counter/internal/api/grpcclient"
	"counter/pkg/api/counter"
)

type AddHandler interface {
	AddHandle(path string, handler http.Handler)
}

type Service struct {
	counter.UnimplementedCounterServer
	listener   net.Listener
	server     *grpc.Server
	router     *net.Listener
	counterSrv CounterService
}

type CounterService interface {
	Increase(ctx context.Context, postId int32) error
	Decrease(ctx context.Context, postId int32) error
}

func (s *Service) Increase(ctx context.Context, request *counter.CounterRequest) (*counter.CounterResponse, error) {
	err := s.counterSrv.Increase(ctx, request.PostId)
	if err != nil {
		return &counter.CounterResponse{Status: false}, err
	}
	return &counter.CounterResponse{Status: true}, err
}

func (s *Service) Decrease(ctx context.Context, request *counter.CounterRequest) (*counter.CounterResponse, error) {
	err := s.counterSrv.Decrease(ctx, request.PostId)
	if err != nil {
		return &counter.CounterResponse{Status: false}, err
	}
	return &counter.CounterResponse{Status: true}, err
}

func NewGRPCServer(cfg *grpcclient.ConfigServer, repo CounterService) (*Service, error) {
	srv := new(Service)
	var err error
	srv.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	srv.server = grpc.NewServer()
	srv.counterSrv = repo
	reflection.Register(srv.server)

	return srv, nil
}

func (s *Service) Start() {
	go func() {
		err := s.server.Serve(s.listener)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Service) Register(ctx context.Context, service AddHandler) error {
	counter.RegisterCounterServer(s.server, s)
	mux := runtime.NewServeMux()
	if err := counter.RegisterCounterHandlerServer(ctx, mux, s); err != nil {
		return err
	}
	service.AddHandle("/", mux)
	return nil
}
