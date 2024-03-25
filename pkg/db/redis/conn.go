package redis

import (
	"github.com/redis/go-redis/v9"
	"task/config"
	"time"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	redisHost := cfg.Redis.Addr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  cfg.Redis.PoolTimeout * time.Second,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
	})

	return client
}
