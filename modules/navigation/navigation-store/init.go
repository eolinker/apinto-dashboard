package navigation_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		navigation := newNavigationStore(db)
		navigationModule := newNavigationStore(db)
		bean.Injection(&navigation)
		bean.Injection(&navigationModule)
	})
}
