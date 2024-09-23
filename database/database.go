package database

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func Init() *redis.Client {
	ctx := context.Background()
	// Initialize database connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		panic(err)
	}

	slog.Info("Connected to Redis")

	return rdb
}
