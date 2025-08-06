package db

import (
	"context"

	"github.com/gabrielluizsf/rinha-backend-2005/env"
	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr: env.Get().Redis,
	})
}
