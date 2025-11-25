package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/kafka"
	"github.com/amogh-sood/distributed-autoscaler/api/internal/metrics"
	"github.com/amogh-sood/distributed-autoscaler/api/internal/redis"
)

type HashRequest struct {
	Input string `json:"input"`
}

type HashHandler struct {
	producer *kafka.Producer
	redis    *redis.Client
}

func NewHashHandler(producer *kafka.Producer, redisClient *redis.Client) *HashHandler {
	return &HashHandler{
		producer: producer,
		redis:    redisClient,
	}
}

func (h *HashHandler) HandleHash(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	metrics.APIRequests.WithLabelValues("/hash", "POST").Inc()

	var req HashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.APIErrors.WithLabelValues("/hash", "POST", "400").Inc()
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	jobID := uuid.NewString()
	payload := fmt.Sprintf(`{"job_id":"%s","input":"%s"}`, jobID, req.Input)

	enqueueStart := time.Now()
	if err := h.producer.Publish(payload); err != nil {
		metrics.APIErrors.WithLabelValues("/hash", "POST", "500").Inc()
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}
	metrics.APIKafkaEnqueueDuration.Observe(time.Since(enqueueStart).Seconds())

	metrics.APIRequestDuration.WithLabelValues("/hash", "POST").
		Observe(time.Since(start).Seconds())

	// Return job ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"job_id":"%s"}`, jobID)))
}
