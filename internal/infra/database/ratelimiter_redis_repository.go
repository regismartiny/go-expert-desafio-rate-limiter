package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RateLimiterRedisRepository struct {
	ctx    context.Context
	client *redis.Client
}

func NewRateLimiterRedisRepository(ctx context.Context, client *redis.Client) *RateLimiterRedisRepository {
	return &RateLimiterRedisRepository{ctx: ctx, client: client}
}

func (r *RateLimiterRedisRepository) SaveActiveClients(clients map[string]ActiveClient) error {

	value, err := json.Marshal(clients)
	if err != nil {
		return err
	}

	for _, client := range clients {
		err = r.client.Set(r.ctx, client.ClientId, value, 0).Err()
		if err != nil {
			fmt.Println("Error saving active client to Redis", err)
		}
	}

	return err
}

func (r *RateLimiterRedisRepository) GetActiveClients() (map[string]ActiveClient, error) {

	activeClients := make(map[string]ActiveClient, 10)

	iter := r.client.Scan(r.ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(r.ctx) {
		val := iter.Val()
		fmt.Println("keys", val)

		var activeClient ActiveClient
		err := json.Unmarshal([]byte(val), &activeClient)
		if err != nil {
			return map[string]ActiveClient{}, err
		}

		activeClients[activeClient.ClientId] = activeClient
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return activeClients, nil
}
