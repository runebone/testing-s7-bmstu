package cache

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(addr string, ttl int) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &Cache{
		client: rdb,
		ttl:    time.Duration(ttl) * time.Second,
	}
}

func (c *Cache) GetOrSet(ctx context.Context, key string, fetchFunc func() (interface{}, error)) (interface{}, error) {
	cachedData, err := c.client.Get(ctx, key).Result()

	if err == redis.Nil {
		data, err := fetchFunc()

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.Marshal(data)
		c.client.Set(ctx, key, jsonData, c.ttl)

		return data, nil
	} else if err != nil {
		return nil, err
	}

	var result interface{}
	json.Unmarshal([]byte(cachedData), &result)

	return result, nil
}
