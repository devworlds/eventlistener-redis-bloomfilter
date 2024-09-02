package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/devworlds/eventlistener-redis-performance/internal/db"
	"github.com/devworlds/eventlistener-redis-performance/internal/listener"
	"github.com/devworlds/eventlistener-redis-performance/internal/redis"
)

func listenInit() {
	ethClient := listener.Connect()
	rdb := redis.Connect()
	_ = rdb
	db.Connect()
	filter := bloom.NewWithEstimates(1000000, 0.01)
	_ = filter

	pubsub := rdb.SubscribeToChan()
	ch := pubsub.Channel()
	fmt.Println("Listening Channel...")

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	go listener.Listen(ethClient, filter, rdb, ctx)
	for msg := range ch {
		fmt.Printf("New message: %s\n", msg.Payload)
		cancel()

		payload, err := rdb.GetWallet()
		if err != nil {
			log.Fatal(err)
		}
		filter.UnmarshalJSON([]byte(payload))
		fmt.Printf("restarting listener..\n")
		go listener.Listen(ethClient, filter, rdb, ctx)
	}
}
