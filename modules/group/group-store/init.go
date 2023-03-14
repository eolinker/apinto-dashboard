package group_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		commonGroup := newCommonGroupStore(db)

		bean.Injection(&commonGroup)
	})
}
