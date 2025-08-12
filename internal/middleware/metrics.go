package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Métricas
var (
	// Contador para o número total de requisições HTTP
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Número total de requisições HTTP",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Histograma para a duração das requisições HTTP
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP",
			Buckets: prometheus.DefBuckets, // Buckets padrão
		},
		[]string{"method", "endpoint"},
	)
)

// MetricsMiddleware é um middleware que coleta métricas de requisições HTTP.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrapper para capturar o código de status
		rw := &responseWriterDelegator{ResponseWriter: w, status: http.StatusOK}

		// Chama o próximo handler
		next.ServeHTTP(rw, r)

		// Calcula a duração
		duration := time.Since(start).Seconds()

		// Atualiza as métricas
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(rw.status)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

// responseWriterDelegator é um wrapper para http.ResponseWriter que captura o código de status.
type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	if r.written {
		return // Evita múltiplas chamadas a WriteHeader
	}
	r.status = code
	r.written = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.written {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// MetricsHandler retorna um handler para o endpoint /metrics do Prometheus.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}