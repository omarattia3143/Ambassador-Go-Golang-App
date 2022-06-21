package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

var Cache *redis.Client
var cacheChannel chan string

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("connected to redis")
}

func SetupCacheChannel() {
	cacheChannel = make(chan string)

	go func(ch chan string) {
		for {
			key := <-ch
			Cache.Del(context.Background(), key)
			println(key + " Cleared!")
		}
	}(cacheChannel)
}

func ClearCache(keys ...string) {

	for _, key := range keys {
		cacheChannel <- key
	}
}