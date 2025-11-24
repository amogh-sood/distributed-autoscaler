package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	APIRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Number of API requests received",
		},
		[]string{"endpoint", "method"},
	)

	APIErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_errors_total",
			Help: "Number of failed API requests",
		},
		[]string{"endpoint", "method", "code"},
	)

	APIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Request duration for API endpoints",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	APIKafkaEnqueueDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "api_kafka_enqueue_seconds",
			Help:    "Time taken to enqueue a job into Kafka",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	// Prometheus auto-registers promauto metrics.
	// This function exists for symmetry and future expansion.
}
