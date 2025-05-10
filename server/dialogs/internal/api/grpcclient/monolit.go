package grpcclient

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "dialogs/pkg/api"
)

type MonolitService struct {
	client pb.MonolitClient
}

func NewMonolitService(cfg Config) *MonolitService {
	conn, err := createConnection(cfg)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &MonolitService{
		client: pb.NewMonolitClient(conn),
	}
}

func (receiver MonolitService) HealthCheck(ctx context.Context) bool {
	e := emptypb.Empty{}
	result, err := receiver.client.HealthCheck(ctx, &e)
	if err != nil {
		return false
	}
	return result.Status == "ok"
}
