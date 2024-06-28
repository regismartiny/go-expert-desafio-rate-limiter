package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RateLimiterRedisRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewRateLimiterRedisRepository(ctx context.Context, client *redis.Client) *RateLimiterRedisRepository {
	return &RateLimiterRedisRepository{ctx: ctx, client: client}
}

func (r *RateLimiterRedisRepository) SaveRequestCount(key string, value int) error {

	err := r.client.Set(r.ctx, key, value, 0).Err()

	if err != nil {
		fmt.Println("Error saving reqCount to Redis", err)
	}

	return err
}

func (r *RateLimiterRedisRepository) GetRequestCount(key string) (int, error) {

	reqCountStr, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting reqCount from Redis", err)
	}

	reqCount, err := strconv.Atoi(reqCountStr)
	if err != nil {
		fmt.Println("Error converting reqCount to integer", err)
		return 0, err
	}

	return reqCount, nil
}
