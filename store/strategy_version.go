package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/strategy-entry"
	"github.com/eolinker/apinto-dashboard/entry/version-entry"
)

type IStrategyVersionStore interface {
	IBaseStore[strategy_entry.StrategyVersion]
}

type strategyVersionHandler struct {
}

func (s *strategyVersionHandler) Kind() string {
	return "strategy"
}

func (s *strategyVersionHandler) Encode(sv *strategy_entry.StrategyVersion) *version_entry.Version {

	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.StrategyId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId

	bytes, _ := json.Marshal(sv.StrategyConfigInfo)
	v.Data = bytes

	return v
}

func (s *strategyVersionHandler) Decode(v *version_entry.Version) *strategy_entry.StrategyVersion {
	sv := new(strategy_entry.StrategyVersion)
	sv.Id = v.Id
	sv.StrategyId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime

	_ = json.Unmarshal(v.Data, &sv.StrategyConfigInfo)

	return sv
}

func newStrategyVersionStore(db IDB) IStrategyVersionStore {
	var h BaseKindHandler[strategy_entry.StrategyVersion, version_entry.Version] = &strategyVersionHandler{}
	return CreateBaseKindStore(h, db)
}
