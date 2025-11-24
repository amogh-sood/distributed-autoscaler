package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var storeCtx = context.Background()

type Store struct {
	rdb *redis.Client
}

func NewStore(addr string) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Store{rdb: rdb}
}

func (s *Store) WriteResult(jobID string, data string) error {
	return s.rdb.Set(storeCtx, "job:"+jobID, data, 30*time.Minute).Err()
}
