package main

import (
	"log"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/server"
)

func main() {
	srv := server.New()

	log.Println("Starting API service...")
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("API failed: %v", err)
	}
}
