package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status bool   `json:"status"`
	Time   string `json:"time"`
}

type Client interface {
	HealthCheck(ctx context.Context) bool
}

type HealthCheckHandler struct {
	uptime time.Time
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{
		uptime: time.Now(),
	}
}

func (hc *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := HealthResponse{
		Status: true,
		Time:   fmt.Sprintf("%.4f", time.Now().Sub(hc.uptime).Hours()),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
