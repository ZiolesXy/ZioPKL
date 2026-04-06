package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() *redis.Client {
	rh := os.Getenv("REDIS_HOST")
	rp := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", rh, rp),
		Password: "",
		DB: 0,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	RDB = rdb
	log.Println("Redis succesfully connected")
	return rdb
}
