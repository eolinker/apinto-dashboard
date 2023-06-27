/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Locker interface {
	TryLock(ctx context.Context, key string, expire time.Duration) (bool, error)
	Unlock(ctx context.Context, key string)
}

type redisLocker struct {
	client *redis.ClusterClient
}

func (r *redisLocker) TryLock(ctx context.Context, key string, expire time.Duration) (bool, error) {
	cmd := r.client.SetNX(ctx, key, 1, expire)
	err := cmd.Err()
	if err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

func (r *redisLocker) Unlock(ctx context.Context, key string) {
	//TODO implement me
	panic("implement me")
}
