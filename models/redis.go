package models

import (
	"eda/logger"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"os"
)

var RedisClient *redis.Client

func NewRedis() error {
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		addr = "redis:6379"
	}

	password, _ := os.LookupEnv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logger.Log.Error("Failed to connect to redis", zap.Error(err))
		return err
	}

	RedisClient = client
	return nil
}

func GetDelPhoneTransaction(phone string) (string, error) {
	pipe := RedisClient.TxPipeline()

	currentCodeCmd := pipe.Get(phone)
	pipe.Del(phone)

	_, err := pipe.Exec()
	if err != nil {
		return "", err
	}

	code, err := currentCodeCmd.Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no code found for phone: %s", phone)
	} else if err != nil {
		return "", err
	}

	return code, nil
}
