package middleware

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds.",
			// Default buckets are ok; adapt if you expect very different latencies.
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	// only call registration once per process
	once sync.Once
)

// RegisterMetrics registers the metrics to the default registry.
// It is safe to call multiple times; registration only happens once.
// If the collector is already registered, we reuse the existing collector
// instead of panicking.
func RegisterMetrics() {
	once.Do(func() {
		register := func(c prometheus.Collector) {
			if err := prometheus.Register(c); err != nil {
				if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
					// reuse the already-registered collector
					switch existing := are.ExistingCollector.(type) {
					case *prometheus.HistogramVec:
						if c == RequestDuration {
							RequestDuration = existing
						}
					case *prometheus.CounterVec:
						if c == RequestCounter {
							RequestCounter = existing
						}
					default:
						// ignore other types
					}
				} else {
					log.Fatalf("failed to register collector: %v", err)
				}
			}
		}

		register(RequestDuration)
		register(RequestCounter)
	})
}

// MetricsHandler returns an http.Handler that serves metrics and ensures
// RegisterMetrics() is called first.
func MetricsHandler() http.Handler {
	RegisterMetrics()
	return promhttp.Handler()
}
