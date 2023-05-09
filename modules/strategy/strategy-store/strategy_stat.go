package strategy_store

import (
	"github.com/eolinker/apinto-dashboard/modules/base/stat-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IStrategyStatStore interface {
	store.IBaseStore[strategy_entry.StrategyStat]
}

type strategyStatKindHandler struct {
}

func (s *strategyStatKindHandler) Kind() string {
	return "strategy"
}

func (s *strategyStatKindHandler) Encode(sv *strategy_entry.StrategyStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.StrategyId
	stat.Kind = s.Kind()
	stat.Version = sv.VersionId

	return stat
}

func (s *strategyStatKindHandler) Decode(stat *stat_entry.Stat) *strategy_entry.StrategyStat {
	ds := new(strategy_entry.StrategyStat)

	ds.StrategyId = stat.Tag
	ds.VersionId = stat.Version

	return ds
}

func newStrategyStatStore(db store.IDB) IStrategyStatStore {
	var h store.BaseKindHandler[strategy_entry.StrategyStat, stat_entry.Stat] = new(strategyStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
