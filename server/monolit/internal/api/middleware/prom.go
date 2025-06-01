package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HandleProm struct {
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func NewHandleProm() *HandleProm {
	result := &HandleProm{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{0.001, 0.01, 0.10, .30, .50, .80, 1},
			},
			[]string{"method", "path"},
		),
	}
	prometheus.MustRegister(result.httpRequestsTotal, result.httpRequestDuration)
	return result
}

func (hp *HandleProm) PrometheusMiddleware(next HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку для захвата статус-кода
		rw := &responseWriter{w, http.StatusOK}
		defer func() {
			duration := time.Since(start).Seconds()

			// Метки для метрик
			method := r.Method
			path := normalizePath(r.URL.Path) // Нормализуем путь (см. ниже)
			status := http.StatusText(rw.status)

			// Записываем метрики
			hp.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
			hp.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
		}()

		next(rw, r)
	}
}

// Обертка для захвата статус-кода
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Нормализация пути (например, /user/123 → /user/:id)
func normalizePath(path string) string {
	// Реализуйте логику замены динамических сегментов
	return path
}
