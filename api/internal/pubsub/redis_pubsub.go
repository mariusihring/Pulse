package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub(addr, password string, db int) *RedisPubSub {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisPubSub{client: client}
}

func (rps *RedisPubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rps.client.Publish(ctx, channel, data).Err()
}

func (rps *RedisPubSub) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return rps.client.Subscribe(ctx, channel)
}
