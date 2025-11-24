package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JobDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "worker_job_duration_seconds",
		Help:    "Time spent processing a single job",
		Buckets: prometheus.DefBuckets,
	})

	NonceAttempts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "worker_nonce_attempts_total",
		Help: "Total number of nonce attempts across all jobs",
	})

	JobsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "worker_jobs_processed_total",
		Help: "Total number of completed jobs",
	})
)

func Register() {
	prometheus.MustRegister(JobDuration, NonceAttempts, JobsProcessed)
}
