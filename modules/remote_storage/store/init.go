package store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		remoteStorageStore := newRemoteStorageStore(db)
		bean.Injection(&remoteStorageStore)
	})
}
