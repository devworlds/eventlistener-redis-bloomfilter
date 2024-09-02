package redis

import "github.com/redis/go-redis/v9"

type RedisService struct {
	client *redis.Client
}
