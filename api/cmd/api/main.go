package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/server"
)

func main() {
	log.Println("Starting API service...")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("API metrics on :9090/metrics")
		http.ListenAndServe(":9090", nil)
	}()

	srv := server.New()
	srv.Start(":8080")
}
