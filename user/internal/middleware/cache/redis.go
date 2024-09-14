package cache

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type CacheMiddleware struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewCacheMiddleware(addr string, ttl int) *CacheMiddleware {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &CacheMiddleware{
		redisClient: rdb,
		ttl:         time.Duration(ttl) * time.Second,
	}
}

func (c *CacheMiddleware) GetOrSet(ctx context.Context, key string, fetchFunc func() (interface{}, error)) (interface{}, error) {
	cachedData, err := c.redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		data, err := fetchFunc()

		if err != nil {
			return nil, err
		}

		jsonData, _ := json.Marshal(data)
		c.redisClient.Set(ctx, key, jsonData, c.ttl)

		return data, nil
	} else if err != nil {
		return nil, err
	}

	var result interface{}
	json.Unmarshal([]byte(cachedData), &result)

	return result, nil
}
