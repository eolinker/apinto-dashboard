package dynamic_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

var (
	kind = "dynamic_module"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		dynamic := newDynamicStore(db)
		dynamicPublishHistory := NewDynamicPublishHistoryStore(db, kind)
		dynamicPublishVersion := NewDynamicPublishVersionStore(db, kind)
		bean.Injection(&dynamic)
		bean.Injection(&dynamicPublishHistory)
		bean.Injection(&dynamicPublishVersion)
	})
}
