package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/amogh-sood/distributed-autoscaler/worker/internal/consumer"
	"github.com/amogh-sood/distributed-autoscaler/worker/internal/metrics"
	"github.com/amogh-sood/distributed-autoscaler/worker/internal/redis"
)

func main() {
	log.Println("Worker: starting up...")

	metrics.Register()

	// Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Worker metrics on :9090/metrics")
		http.ListenAndServe(":9090", nil)
	}()

	// Redis
	redisClient := redis.New("localhost:6379")

	// Kafka consumer
	c := consumer.NewConsumer("localhost:9092", "jobs", "worker-group", redisClient)
	c.Start()
}
