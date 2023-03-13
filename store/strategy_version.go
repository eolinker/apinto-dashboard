package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IStrategyVersionStore interface {
	IBaseStore[entry.StrategyVersion]
}

type strategyVersionHandler struct {
}

func (s *strategyVersionHandler) Kind() string {
	return "strategy"
}

func (s *strategyVersionHandler) Encode(sv *entry.StrategyVersion) *entry.Version {

	v := new(entry.Version)
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

func (s *strategyVersionHandler) Decode(v *entry.Version) *entry.StrategyVersion {
	sv := new(entry.StrategyVersion)
	sv.Id = v.Id
	sv.StrategyId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime

	_ = json.Unmarshal(v.Data, &sv.StrategyConfigInfo)

	return sv
}

func newStrategyVersionStore(db IDB) IStrategyVersionStore {
	var h BaseKindHandler[entry.StrategyVersion, entry.Version] = &strategyVersionHandler{}
	return CreateBaseKindStore(h, db)
}
