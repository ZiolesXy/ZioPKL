package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: "",
		DB: 0,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	RDB = rdb
	return rdb
}
