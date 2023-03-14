package locker_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	lockAsynService := NewAsynLockService()
	lockSyncService := NewSyncLockService()
	bean.Injection(&lockAsynService)
	bean.Injection(&lockSyncService)
}
