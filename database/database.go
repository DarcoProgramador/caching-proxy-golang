package database

import (
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Init() {
	// Initialize database connection
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetClient() *redis.Client {
	return rdb
}
