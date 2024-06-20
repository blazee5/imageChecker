package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
}

func NewRedisClient(cfg Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		panic(err)
	}

	return client
}
