package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status bool   `json:"status"`
	Time   string `json:"time"`
}

func (srv *Service) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := HealthResponse{
		Status: true,
		Time:   fmt.Sprintf("%.4f", time.Now().Sub(srv.uptime).Hours()),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
