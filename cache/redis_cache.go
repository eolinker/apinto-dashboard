package cache

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc/common/bean"
	"strings"
	"time"
)

type IRedisCache[T any, K comparable] interface {
	Get(ctx context.Context, k K) (*T, error)
	Set(ctx context.Context, k K, t *T) error

	Delete(ctx context.Context, keys ...K) error
}

type IRedisCacheSingleton[T any] interface {
	Get(ctx context.Context) (*T, error)
	Set(ctx context.Context, t *T) error
	Delete(ctx context.Context) error
}

type IRedisCacheNoKey[T any] interface {
	SetAll(ctx context.Context, t []*T) error
	GetAll(ctx context.Context) ([]*T, error)
}
type redisCache[T any, K comparable] struct {
	client        ICommonCache
	formatHandler func(K) string
	expiration    time.Duration
}
type redisCacheSingleton[T any] struct {
	base IRedisCache[T, string]
	key  string
}

func (r *redisCacheSingleton[T]) Get(ctx context.Context) (*T, error) {
	return r.base.Get(ctx, r.key)
}

func (r *redisCacheSingleton[T]) Set(ctx context.Context, t *T) error {
	return r.base.Set(ctx, r.key, t)
}

func (r *redisCacheSingleton[T]) Delete(ctx context.Context) error {
	return r.base.Delete(ctx, r.key)
}

type redisCacheNoKey[T any] struct {
	client     ICommonCache
	key        string
	expiration time.Duration
}

func (r *redisCache[T, K]) Get(ctx context.Context, k K) (*T, error) {
	kv := r.formatHandler(k)

	bytes, err := r.client.Get(ctx, kv)
	if err != nil {
		return nil, err
	}

	return r.toMyStruct(bytes)

}

func (r *redisCache[T, K]) Set(ctx context.Context, k K, t *T) error {

	kv := r.formatHandler(k)

	bytes, err := structToBytes(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, kv, bytes, r.expiration)
}

func (r *redisCache[T, K]) Delete(ctx context.Context, ks ...K) error {
	for _, k := range ks {
		key := r.formatHandler(k)
		if err := r.client.Del(ctx, key); err != nil {
			return err
		}

	}
	return nil
}

func (r *redisCacheNoKey[T]) GetAll(ctx context.Context) ([]*T, error) {

	bytes, err := r.client.Get(ctx, r.key)
	if err != nil {
		return nil, err
	}

	return toMyStructAll[T](bytes)

}

func (r *redisCacheNoKey[T]) SetAll(ctx context.Context, t []*T) error {

	bytes, err := structToBytesAll(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.key, bytes, r.expiration)
}
func CreateRedisCacheNoKey[T any](expiration time.Duration, key string, prefix ...string) IRedisCacheNoKey[T] {
	keyPrefix := "apinto-dashboard:"
	if len(key) > 0 {
		keyPrefix = strings.Join(prefix, ":")
	}
	r := &redisCacheNoKey[T]{
		key:        key,
		expiration: expiration,
	}
	var c ICommonCache

	bean.Autowired(&c)
	bean.AddInitializingBeanFunc(func() {
		r.client = c.clone(keyPrefix)
	})
	return r
}

func CreateRedisCache[T any, K comparable](expiration time.Duration, format func(k K) string, key ...string) IRedisCache[T, K] {
	keyPrefix := "apinto-dashboard:"
	if len(key) > 0 {
		keyPrefix = strings.Join(key, "-")
	}
	r := &redisCache[T, K]{
		formatHandler: format,
		expiration:    expiration,
	}
	var c ICommonCache

	bean.Autowired(&c)
	bean.AddInitializingBeanFunc(func() {
		r.client = c.clone(keyPrefix)
	})
	return r
}
func CreateRedisCacheSingleton[T any](expiration time.Duration, key string) IRedisCacheSingleton[T] {
	return &redisCacheSingleton[T]{
		base: CreateRedisCache[T, string](expiration, func(k string) string {
			return k
		}),
		key: key,
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
