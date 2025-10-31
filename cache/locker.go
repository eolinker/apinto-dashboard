package cache

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/log"
	"time"

	"github.com/redis/go-redis/v9"
)

type tRedisLock struct {
	key         string
	value       string
	redisClient redis.UniversalClient
	expiration  time.Duration

	cancelFunc context.CancelFunc
	ctx        context.Context
}

const (
	PubSubPrefix        = "{redis_lock}_"
	DefaultExpiration   = 30
	DefaultSpinInterval = 100
)

func newRedisLock(redisClient redis.UniversalClient, key, value string) Locker {
	return &tRedisLock{
		key:         key,
		value:       value,
		redisClient: redisClient,
		expiration:  time.Duration(DefaultExpiration) * time.Second}
}

func newRedisLockWithExpireTime(redisClient redis.UniversalClient, key string, value string, expiration time.Duration) Locker {
	return &tRedisLock{
		key:         key,
		value:       value,
		redisClient: redisClient,
		expiration:  expiration}
}

// TryLock try get lock only once, if get the lock return true, else return false
func (lock *tRedisLock) TryLock() (bool, error) {
	success, err := lock.redisClient.SetNX(context.Background(), lock.key, lock.value, lock.expiration).Result()
	if err != nil {
		return false, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	lock.cancelFunc = cancelFunc
	lock.renew(ctx)

	return success, nil
}

// Lock blocked until get lock
func (lock *tRedisLock) Lock() error {
	for {
		success, err := lock.TryLock()
		if err != nil {
			return err
		}
		if success {
			return nil
		}
		if !success {
			err := lock.subscribeLock()
			if err != nil {
				return err
			}
		}
	}
}

// Unlock release the lock
func (lock *tRedisLock) Unlock() error {
	script := redis.NewScript(fmt.Sprintf(
		`if redis.call("get", KEYS[1]) == "%s" then return redis.call("del", KEYS[1]) else return 0 end`,
		lock.value))
	if lock.cancelFunc != nil {
		lock.cancelFunc() //cancel renew goroutine
		lock.cancelFunc = nil
	}
	runCmd := script.Run(context.Background(), lock.redisClient, []string{lock.key})
	res, err := runCmd.Result()
	if err != nil {
		return err
	}
	if tmp, ok := res.(int64); ok {
		if tmp == 1 {

			err := lock.publishLock()
			if err != nil {
				return err
			}
			return nil
		}
	}
	err = fmt.Errorf("unlock script fail: %s", lock.key)
	return err
}

// LockWithTimeout blocked until get lock or timeout
func (lock *tRedisLock) LockWithTimeout(d time.Duration) error {
	timeNow := time.Now()
	for {
		success, err := lock.TryLock()
		if err != nil {
			return err
		}
		if success {
			return nil
		}
		deltaTime := d - time.Since(timeNow)
		if !success {
			err := lock.subscribeLockWithTimeout(deltaTime)
			if err != nil {
				return err
			}
		}
	}
}

func (lock *tRedisLock) SpinLock(times int) error {
	for i := 0; i < times; i++ {
		success, err := lock.TryLock()
		if err != nil {
			return err
		}
		if success {
			return nil
		}
		time.Sleep(time.Millisecond * DefaultSpinInterval)
	}
	return fmt.Errorf("max spin times reached")
}

// subscribeLock blocked until lock is released
func (lock *tRedisLock) subscribeLock() error {
	pubSub := lock.redisClient.Subscribe(context.Background(), getPubSubTopic(lock.key))
	_, err := pubSub.Receive(context.Background())
	if err != nil {
		return err
	}
	<-pubSub.Channel()
	return nil
}

// subscribeLock blocked until lock is released or timeout
func (lock *tRedisLock) subscribeLockWithTimeout(d time.Duration) error {
	timeNow := time.Now()
	pubSub := lock.redisClient.Subscribe(context.Background(), getPubSubTopic(lock.key))
	_, err := pubSub.ReceiveTimeout(context.Background(), d)
	if err != nil {
		return err
	}
	deltaTime := time.Since(timeNow) - d
	select {
	case <-pubSub.Channel():
		return nil
	case <-time.After(deltaTime):
		return fmt.Errorf("timeout")
	}
}

// publishLock publish a message about lock is released
func (lock *tRedisLock) publishLock() error {
	err := lock.redisClient.Publish(context.Background(), getPubSubTopic(lock.key), "release lock").Err()
	if err != nil {
		return err
	}
	return nil
}

// renew renew the expiration of lock, and can be canceled when call Unlock
func (lock *tRedisLock) renew(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(lock.expiration / 3)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Debug("reset lok expire:", lock.key)
				lock.redisClient.Expire(context.Background(), lock.key, lock.expiration).Result()
			}
		}
	}()
}

// getPubSubTopic key -> PubSubPrefix + key
func getPubSubTopic(key string) string {
	return PubSubPrefix + key
}
