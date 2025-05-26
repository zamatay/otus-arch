package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/zamatay/otus/arch/lesson-1/internal/api/grpcclient/counter"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/middleware"
	"github.com/zamatay/otus/arch/lesson-1/internal/api/utils"
	"github.com/zamatay/otus/arch/lesson-1/internal/grpcserver"
)

type AddRouted interface {
	AddRoute(string, func(http.ResponseWriter, *http.Request))
	AddProtectedRoute(string, func(http.ResponseWriter, *http.Request))
	AddHandle(string, http.Handler)
}

type Config struct {
	Host         string
	Port         uint16
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Service struct {
	http.Server
	router *http.ServeMux
	//protectedRouter *http.ServeMux
	ErrorCh    chan error
	uptime     time.Time
	Grpc       *grpcserver.Service
	CounterSrv *counter.CounterService
	handleProm *middleware.HandleProm
}

func New(config *Config, secret string) (*Service, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	srv := new(Service)
	srv.router = http.NewServeMux()
	//srv.protectedRouter = http.NewServeMux()

	srv.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	srv.Handler = srv.router
	srv.ReadTimeout = config.ReadTimeout
	srv.WriteTimeout = config.WriteTimeout
	srv.IdleTimeout = config.IdleTimeout

	utils.UserByToken = utils.SetUserByToken([]byte(secret))

	srv.router.HandleFunc("/health", srv.healthCheckHandler)
	srv.router.Handle("/metrics", promhttp.Handler())

	srv.uptime = time.Now()
	srv.ErrorCh = make(chan error)

	srv.handleProm = middleware.NewHandleProm()
	return srv, nil
}

func (srv *Service) AddRoute(path string, handler func(http.ResponseWriter, *http.Request)) {
	srv.router.HandleFunc(path, middleware.CorsMiddleware(handler))
}

func (srv *Service) AddProtectedRoute(path string, handler func(http.ResponseWriter, *http.Request)) {
	srv.router.HandleFunc(path, srv.handleProm.PrometheusMiddleware(middleware.CorsMiddleware(middleware.TokenMiddleware((handler)))))
}

func (srv *Service) AddHandle(path string, handler http.Handler) {
	srv.router.Handle(path, handler)
}

func (srv *Service) Start() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	go func() {
		err := srv.Serve(ln)
		if err != nil {
			srv.ErrorCh <- err
		}
	}()

	go func() {
		if err := srv.Grpc.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (srv *Service) Stop(ctx context.Context) error {
	srv.Server.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (srv *Service) GetRoute() *http.ServeMux {
	return srv.router
}
