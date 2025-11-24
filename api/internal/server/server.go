package server

import (
	"log"
	"net/http"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/handlers"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/kafka"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router   *chi.Mux
	producer *kafka.Producer
}

func New() *Server {
	producer := kafka.NewProducer("localhost:9092", "jobs")

	r := chi.NewRouter()

	h := handlers.NewHashHandler(producer)
	r.Post("/hash", h.HandleHash)

	return &Server{
		router:   r,
		producer: producer,
	}
}

func (s *Server) Start(addr string) error {
	log.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
