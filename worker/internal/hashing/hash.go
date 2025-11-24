package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/amogh-sood/distributed-autoscaler/worker/internal/metrics"
)

func Mine(input string) (string, int) {
	start := time.Now()

	nonce := 0
	var hash string

	for {
		nonce++
		data := fmt.Sprintf("%s:%d", input, nonce)
		sum := sha256.Sum256([]byte(data))
		hash = hex.EncodeToString(sum[:])

		metrics.NonceAttempts.Inc()

		// found a hash with leading zeros
		if hash[:4] == "0000" {
			break
		}
	}

	duration := time.Since(start).Seconds()
	metrics.JobDuration.Observe(duration)
	metrics.JobsProcessed.Inc()

	return hash, nonce
}
