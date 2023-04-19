package store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		warnStrategy := newWarnStrategyIStore(db)
		warnHistory := newWarnHistoryIStore(db)
		bean.Injection(&warnStrategy)
		bean.Injection(&warnHistory)
	})
}
