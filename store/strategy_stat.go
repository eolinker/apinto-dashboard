package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IStrategyStatStore interface {
	IBaseStore[entry.StrategyStat]
}

type strategyStatKindHandler struct {
}

func (s *strategyStatKindHandler) Kind() string {
	return "strategy"
}

func (s *strategyStatKindHandler) Encode(sv *entry.StrategyStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.StrategyId
	stat.Kind = s.Kind()
	stat.Version = sv.VersionId

	return stat
}

func (s *strategyStatKindHandler) Decode(stat *entry.Stat) *entry.StrategyStat {
	ds := new(entry.StrategyStat)

	ds.StrategyId = stat.Tag
	ds.VersionId = stat.Version

	return ds
}

func newStrategyStatStore(db IDB) IStrategyStatStore {
	var h BaseKindHandler[entry.StrategyStat, entry.Stat] = new(strategyStatKindHandler)
	return CreateBaseKindStore(h, db)
}
