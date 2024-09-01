package main

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/devworlds/eventlistener-redis-performance/internal/db"
	"github.com/devworlds/eventlistener-redis-performance/internal/listener"
	"github.com/devworlds/eventlistener-redis-performance/internal/redis"
)

func main() {
	ethClient := listener.Connect()
	redisClient := redis.Connect()
	dbClient := db.Connect()

	filter := bloom.NewWithEstimates(1000000, 0.01)
	_ = filter

	listener.Listen(ethClient, redisClient, dbClient, filter)
}
