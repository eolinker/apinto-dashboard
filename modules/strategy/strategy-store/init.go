package strategy_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		strategy := newStrategyStore(db)
		strategyStat := newStrategyStatStore(db)

		strategyVersion := newStrategyVersionStore(db)

		strategyHistory := newStrategyHistoryStore(db)
		bean.Injection(&strategy)
		bean.Injection(&strategyStat)

		bean.Injection(&strategyVersion)

		bean.Injection(&strategyHistory)
	})
}
