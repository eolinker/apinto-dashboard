package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type IRedisCache[T any] interface {
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, t *T, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}

type redisCache[T any] struct {
	client    *redis.ClusterClient
	keyPrefix string
}

func (r *redisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	key = r.keyPrefix + key

	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return r.toMyStruct(bytes)

}

func (r *redisCache[T]) Set(ctx context.Context, key string, t *T, expiration time.Duration) error {

	key = r.keyPrefix + key

	bytes, err := r.structToBytes(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, expiration).Err()
}

func (r *redisCache[T]) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		key = r.keyPrefix + key
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}

	}
	return nil
}

func createRedisCache[T any](client *redis.ClusterClient) *redisCache[T] {
	return &redisCache[T]{
		client:    client,
		keyPrefix: "apinto-dashboard:",
	}
}

func (r *redisCache[T]) structToBytes(t *T) ([]byte, error) {

	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func (r *redisCache[T]) toMyStruct(bytes []byte) (*T, error) {

	t := new(T)
	err := json.Unmarshal(bytes, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
