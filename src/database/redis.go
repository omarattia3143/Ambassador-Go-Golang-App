package database

import (
	"fmt"
	"github.com/go-redis/redis/v9"
)

var Cache *redis.Client

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("connected to redis")
}
