package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/handlers"
	"github.com/amogh-sood/distributed-autoscaler/api/internal/kafka"
	"github.com/amogh-sood/distributed-autoscaler/api/internal/redis"
)

type Server struct {
	router   *chi.Mux
	producer *kafka.Producer
	redis    *redis.Client
}

func New() *Server {
	// init kafka
	producer := kafka.NewProducer("localhost:9092", "jobs")

	// init redis
	redisClient := redis.New("localhost:6379")

	r := chi.NewRouter()

	// handler now expects producer + redis
	h := handlers.NewHashHandler(producer, redisClient)
	r.Post("/hash", h.HandleHash)

	return &Server{
		router:   r,
		producer: producer,
		redis:    redisClient,
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("API listening on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
