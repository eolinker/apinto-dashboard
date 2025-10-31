/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package cache

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	lockerPrefix = "locker"
)

type Locker interface {
	TryLock() (bool, error)
	Unlock() error
	Lock() error
}
type RedisLocker interface {
	Locker(name string, value ...string) Locker
	LockerWithExpireTime(name string, expire time.Duration, value ...string) Locker
}
type manager struct {
	client    redis.UniversalClient
	namespace string
}

func newManager(client redis.UniversalClient, namespace string) RedisLocker {
	return &manager{client: client, namespace: namespace}
}

func (m *manager) Locker(name string, value ...string) Locker {
	v := strconv.FormatInt(time.Now().UnixNano(), 10)
	if len(value) > 0 {
		v = strings.Join(value, "-")
	}
	return newRedisLock(client, fmt.Sprintf("%s:%s:%s", m.namespace, lockerPrefix, name), v)

}

func (m *manager) LockerWithExpireTime(name string, expire time.Duration, value ...string) Locker {
	v := strconv.FormatInt(time.Now().UnixNano(), 10)
	if len(value) > 0 {
		v = strings.Join(value, "-")
	}
	return newRedisLockWithExpireTime(client, fmt.Sprintf("%s:%s:%s", m.namespace, lockerPrefix, name), v, expire)

}
