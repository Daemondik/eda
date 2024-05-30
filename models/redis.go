package models

import (
	"errors"
	"github.com/go-redis/redis"
	"os"
)

var RedisClient redis.Client

func NewRedis() {
	var client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if client == nil {
		errors.New("cannot run redis")
	}

	RedisClient = *client
}
