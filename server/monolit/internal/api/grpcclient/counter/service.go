package counter

import (
	"log"

	"github.com/zamatay/otus/arch/lesson-1/internal/api/grpcclient"
	"github.com/zamatay/otus/arch/lesson-1/pkg/api/counter"
)

type CounterService struct {
	Client counter.CounterClient
}

func NewCounterService(cfg grpcclient.Config) *CounterService {
	conn, err := grpcclient.CreateConnection(cfg)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &CounterService{
		Client: counter.NewCounterClient(conn),
	}
}
