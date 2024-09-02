package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func Connect() *RedisService {
	opt, err := redis.ParseURL("redis://:eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81@localhost:6379")
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	fmt.Println("Redis Successfully connected!")

	pong, err := rdb.Ping(context.Background()).Result()
	_ = pong
	if err != nil {
		log.Fatal(err.Error())
	}

	exist := rdb.Exists(context.Background(), "wallets")
	if exist.Val() == 0 {
		rdb.Set(context.Background(), "wallets", "", 0)
	}

	return &RedisService{
		client: rdb,
	}
}

func (r *RedisService) SetWallet(value string) error {
	err := r.client.Set(context.Background(), "wallets", value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetWallet() (string, error) {
	val, err := r.client.Get(context.Background(), "wallets").Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisService) SubscribeToChan() *redis.PubSub {
	sub := r.client.Subscribe(context.Background(), "newWallet")
	_, err := sub.Receive(context.Background())
	if err != nil {
		log.Fatalf("fail to sub chaannel: %v", err)
	}
	return sub
}

func (r *RedisService) PublishToChan() {
	r.client.Publish(context.Background(), "newWallet", "New wallet added.")
}
