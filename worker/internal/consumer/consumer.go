package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/amogh-sood/distributed-autoscaler/worker/internal/hashing"
	"github.com/amogh-sood/distributed-autoscaler/worker/internal/metrics"
	"github.com/amogh-sood/distributed-autoscaler/worker/internal/redis"
)

type Consumer struct {
	reader *kafkago.Reader
	redis  *redis.Client
}

type JobMessage struct {
	JobID string `json:"job_id"`
	Input string `json:"input"`
}

type JobResult struct {
	Status      string  `json:"status"`
	Input       string  `json:"input"`
	Hash        string  `json:"hash"`
	Nonce       int64   `json:"nonce"`
	DurationSec float64 `json:"duration_sec"`
}

func NewConsumer(brokers, topic, group string, redisClient *redis.Client) *Consumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   topic,
		GroupID: group,
	})

	return &Consumer{
		reader: r,
		redis:  redisClient,
	}
}

func (c *Consumer) Start() {
	log.Println("Worker: listening for jobs...")

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Worker: error reading message: %v\n", err)
			continue
		}

		var job JobMessage
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			log.Printf("Worker: failed to unmarshal job: %v (raw=%s)\n", err, string(msg.Value))
			continue
		}

		start := time.Now()
		hash, nonceAttempts := hashing.Mine(job.Input)
		duration := time.Since(start).Seconds()

		// record metrics
		metrics.WorkerJobDuration.Observe(duration)
		metrics.WorkerJobsProcessed.Inc()
		metrics.WorkerNonceAttempts.Add(float64(nonceAttempts))

		log.Printf("Worker: hashed '%s' → %s (nonce=%d)\n", job.Input, hash, nonceAttempts)

		result := JobResult{
			Status:      "done",
			Input:       job.Input,
			Hash:        hash,
			Nonce:       nonceAttempts,
			DurationSec: duration,
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			log.Printf("Worker: failed to marshal result for job %s: %v\n", job.JobID, err)
			continue
		}

		key := "job:" + job.JobID
		if err := c.redis.Set(key, string(resultJSON)); err != nil {
			log.Printf("Worker: failed to write result to Redis for %s: %v\n", key, err)
			continue
		}
	}
}
