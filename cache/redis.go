package cache

import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // Uses redis:6380 from REDIS_ADDR
	})
	return &Cache{Client: client}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(context.Background(), key, data, expiration).Err()
}

func (c *Cache) Get(key string, dest interface{}) error {
	data, err := c.Client.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}