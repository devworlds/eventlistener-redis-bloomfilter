package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func Connect() *redis.Client {
	opt, err := redis.ParseURL("redis://eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81@localhost:6379")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)

	fmt.Println("Redis Successfully connected!")
	return client
}
