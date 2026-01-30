package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis creates and tests a new Redis client with the given configuration.
func NewRedis(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Connected to Redis 🚀")
	return rdb, nil
}

