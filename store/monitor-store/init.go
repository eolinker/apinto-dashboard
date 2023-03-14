package monitor_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		iWarnStrategyStore := newWarnStrategyIStore(db)
		iWarnHistoryStore := newWarnHistoryIStore(db)
		monitor := newMonitorStore(db)
		bean.Injection(&monitor)
		bean.Injection(&iWarnStrategyStore)
		bean.Injection(&iWarnHistoryStore)
	})
}
