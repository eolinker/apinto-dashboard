package locker_service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type LockNameType = string

const (
	LockNameDiscovery       = "discovery"
	LockNameService         = "service"
	LockNameVariable        = "variable"
	LockNameApplication     = "application"
	LockNameAPI             = "api"
	LockNameStrategy        = "strategy"
	LockNameExtApp          = "ext-app"
	LockNameMonPartition    = "monitor-partition"
	LockNamePluginTemplate  = "plugin_template"
	LockNamePluginNamespace = "plugin_namespace"
	LockNameClusterPlugin   = "cluster_plugin"
)

type IAsynLockService interface {
	Lock(name LockNameType, id int) error
	DeleteLock(name LockNameType, id int)
	Unlock(name LockNameType, id int)
}

type ISyncLockService interface {
	Lock(name LockNameType, id int) error
	DeleteLock(name LockNameType, id int)
	Unlock(name LockNameType, id int)
}

// 异步锁
func NewAsynLockService() IAsynLockService {
	lockService := &asynLock{
		lockMaps: &sync.Map{},
	}
	return lockService
}

// 同步锁
func NewSyncLockService() ISyncLockService {
	lockService := &syncLock{
		lockMaps: &sync.Map{},
	}
	return lockService
}

type syncLock struct {
	lockMaps *sync.Map
}

func (s *syncLock) Lock(name LockNameType, id int) error {
	key := fmt.Sprintf("%s_%d", name, id)
	s.lockMaps.LoadOrStore(key, new(sync.Mutex))
	if v, ok := s.lockMaps.Load(key); ok {
		mutex := v.(*sync.Mutex)
		mutex.Lock()
	}
	return nil
}

func (s *syncLock) DeleteLock(name LockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	s.lockMaps.Delete(key)
}

func (s *syncLock) Unlock(name LockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	if v, ok := s.lockMaps.Load(key); ok {
		mutex := v.(*sync.Mutex)
		mutex.Unlock()
	}
}

// 异步锁
type asynLock struct {
	lockMaps *sync.Map
}

func (a *asynLock) Lock(name LockNameType, id int) error {
	key := fmt.Sprintf("%s_%d", name, id)
	a.lockMaps.LoadOrStore(key, new(int32))
	if v, ok := a.lockMaps.Load(key); ok {
		addr := v.(*int32)
		if *addr > 0 {
			return errors.New("已有其他人操作，请稍后。")
		} else {
			atomic.StoreInt32(addr, 1)
		}
	}
	return nil
}

func (a *asynLock) DeleteLock(name LockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	a.lockMaps.Delete(key)
}

func (a *asynLock) Unlock(name LockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	if v, ok := a.lockMaps.Load(key); ok {
		atomic.StoreInt32(v.(*int32), 0)
	}
}
