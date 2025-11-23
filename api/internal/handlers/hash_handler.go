package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/kafka"
)

type HashRequest struct {
	Input string `json:"input"`
}

type HashHandler struct {
	producer *kafka.Producer
}

func NewHashHandler(producer *kafka.Producer) *HashHandler {
	return &HashHandler{producer: producer}
}

func (h *HashHandler) HandleHash(w http.ResponseWriter, r *http.Request) {
	var req HashRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Publish to Kafka queue
	if err := h.producer.Publish(req.Input); err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "queued"}`))
}
