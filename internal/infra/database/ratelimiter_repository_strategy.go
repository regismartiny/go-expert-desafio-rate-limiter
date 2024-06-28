package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	config "github.com/regismartiny/go-expert-desafio-rate-limiter/configs"
)

func RateLimiterRepositoryStrategy(ctx context.Context, configs config.PersistenceConfigs, clientType string) RateLimiterRepository {
	var repository RateLimiterRepository
	switch clientType {
	case "redis":
		repository = NewRateLimiterRedisRepository(ctx, getRedisClient(ctx, configs))
	default:
		repository = NewRateLimiterRedisRepository(ctx, getRedisClient(ctx, configs))
	}
	return repository
}

func getRedisClient(ctx context.Context, configs config.PersistenceConfigs) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.Redis.Addr,
		Password: configs.Redis.Password,
		DB:       configs.Redis.Db,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)

	return redisClient
}
