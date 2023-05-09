package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type IRedisCache[T any, K comparable] interface {
	Get(ctx context.Context, k K) (*T, error)
	Set(ctx context.Context, k K, t *T, expiration time.Duration) error

	Delete(ctx context.Context, keys ...K) error
}
type IRedisCacheNoKey[T any] interface {
	SetAll(ctx context.Context, t []*T, expiration time.Duration) error
	GetAll(ctx context.Context) ([]*T, error)
}
type redisCache[T any, K comparable] struct {
	client        *redis.ClusterClient
	keyPrefix     string
	formatHandler func(K) string
}
type redisCacheNoKey[T any] struct {
	client        *redis.ClusterClient
	keyPrefix     string
	formatHandler func() string
}

func (r *redisCache[T, K]) Get(ctx context.Context, k K) (*T, error) {
	kv := r.keyPrefix + r.formatHandler(k)

	bytes, err := r.client.Get(ctx, kv).Bytes()
	if err != nil {
		return nil, err
	}

	return r.toMyStruct(bytes)

}

func (r *redisCache[T, K]) Set(ctx context.Context, k K, t *T, expiration time.Duration) error {

	kv := r.keyPrefix + r.formatHandler(k)

	bytes, err := structToBytes(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, kv, bytes, expiration).Err()
}

func (r *redisCache[T, K]) Delete(ctx context.Context, ks ...K) error {
	for _, k := range ks {
		key := r.keyPrefix + r.formatHandler(k)
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}

	}
	return nil
}

func (r *redisCacheNoKey[T]) GetAll(ctx context.Context) ([]*T, error) {
	key := r.keyPrefix + r.formatHandler()

	bytes, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return toMyStructAll[T](bytes)

}

func (r *redisCacheNoKey[T]) SetAll(ctx context.Context, t []*T, expiration time.Duration) error {

	key := r.keyPrefix + r.formatHandler()

	bytes, err := structToBytesAll(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, expiration).Err()
}
func CreateRedisCacheNoKey[T any](client *redis.ClusterClient, format func() string, key ...string) IRedisCacheNoKey[T] {
	keyPrefix := "apinto-dashboard:"
	if len(key) > 0 {
		keyPrefix = strings.Join(key, "-")
	}
	return &redisCacheNoKey[T]{
		client:        client,
		keyPrefix:     keyPrefix,
		formatHandler: format,
	}
}

func CreateRedisCache[T any, K comparable](client *redis.ClusterClient, format func(k K) string, key ...string) IRedisCache[T, K] {
	keyPrefix := "apinto-dashboard:"
	if len(key) > 0 {
		keyPrefix = strings.Join(key, "-")
	}
	return &redisCache[T, K]{
		client:        client,
		keyPrefix:     keyPrefix,
		formatHandler: format,
	}
}

func structToBytes[T any](t *T) ([]byte, error) {

	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func (r *redisCache[T, K]) toMyStruct(bytes []byte) (*T, error) {

	t := new(T)
	err := json.Unmarshal(bytes, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func structToBytesAll[T any](t []*T) ([]byte, error) {

	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func toMyStructAll[T any](bytes []byte) ([]*T, error) {

	t := make([]*T, 0)
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
