package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Client{rdb: rdb}
}

func (c *Client) Set(key, value string) error {
	return c.rdb.Set(context.Background(), key, value, 0).Err()
}

func (c *Client) Get(key string) (string, error) {
	return c.rdb.Get(context.Background(), key).Result()
}
