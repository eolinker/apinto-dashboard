package service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type lockNameType string

var (
	lockNameDiscovery    lockNameType = "discovery"
	lockNameService      lockNameType = "service"
	lockNameVariable     lockNameType = "variable"
	lockNameApplication  lockNameType = "application"
	lockNameAPI          lockNameType = "api"
	lockNameStrategy     lockNameType = "strategy"
	lockNameExtApp       lockNameType = "ext-app"
	lockNameMonPartition lockNameType = "monitor-partition"
)

type IAsynLockService interface {
	lock(name lockNameType, id int) error
	deleteLock(name lockNameType, id int)
	unlock(name lockNameType, id int)
}

type ISyncLockService interface {
	lock(name lockNameType, id int) error
	deleteLock(name lockNameType, id int)
	unlock(name lockNameType, id int)
}

// 异步锁
func newAsynLockService() IAsynLockService {
	lockService := &asynLock{
		lockMaps: &sync.Map{},
	}
	return lockService
}

// 同步锁
func newSyncLockService() ISyncLockService {
	lockService := &syncLock{
		lockMaps: &sync.Map{},
	}
	return lockService
}

type syncLock struct {
	lockMaps *sync.Map
}

func (s *syncLock) lock(name lockNameType, id int) error {
	key := fmt.Sprintf("%s_%d", name, id)
	s.lockMaps.LoadOrStore(key, new(sync.Mutex))
	if v, ok := s.lockMaps.Load(key); ok {
		mutex := v.(*sync.Mutex)
		mutex.Lock()
	}
	return nil
}

func (s *syncLock) deleteLock(name lockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	s.lockMaps.Delete(key)
}

func (s *syncLock) unlock(name lockNameType, id int) {
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

func (a *asynLock) lock(name lockNameType, id int) error {
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

func (a *asynLock) deleteLock(name lockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	a.lockMaps.Delete(key)
}

func (a *asynLock) unlock(name lockNameType, id int) {
	key := fmt.Sprintf("%s_%d", name, id)
	if v, ok := a.lockMaps.Load(key); ok {
		atomic.StoreInt32(v.(*int32), 0)
	}
}
