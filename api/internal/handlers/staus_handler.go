package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/amogh-sood/distributed-autoscaler/api/internal/redis"
)

type StatusHandler struct {
	redis *redis.Client
}

func NewStatusHandler(r *redis.Client) *StatusHandler {
	return &StatusHandler{redis: r}
}

func (h *StatusHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	key := "job:" + jobID
	val, err := h.redis.Get(key)

	if err != nil {
		// Result not found → still pending
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"pending"}`))
		return
	}

	// Redis entry exists → return full job result
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))
}
