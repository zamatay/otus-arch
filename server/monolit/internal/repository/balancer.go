package repository

import (
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RandomBalancer struct {
	//rnd      *rand.Rand
	replicas []*pgxpool.Pool
}

func NewRandomBalancer() *RandomBalancer {
	balancer := &RandomBalancer{}
	//balancer.rnd = rand.New(rand.NewSource(time.Now().Unix()))
	return balancer
}
func (b *RandomBalancer) AddReplica(connection *pgxpool.Pool) {
	b.replicas = append(b.replicas, connection)
}

func (b *RandomBalancer) GetReplica() *pgxpool.Pool {
	index := rand.Intn(len(b.replicas))
	return b.replicas[index]
}

func (b *RandomBalancer) GetAllReplica() []*pgxpool.Pool {
	return b.replicas
}
