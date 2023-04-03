package strategy_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/version-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IStrategyPublishVersionStore interface {
	store.IBaseStore[strategy_entry.StrategyPublishVersion]
}

type strategyPublishVersionHandler struct {
	kind string
}

func (s *strategyPublishVersionHandler) Kind() string {
	return s.kind
}

func (s *strategyPublishVersionHandler) Encode(sv *strategy_entry.StrategyPublishVersion) *version_entry.Version {

	for _, publish := range sv.Publish {
		publish.Strategy.NamespaceId = 0
		publish.Strategy.ClusterId = 0
		publish.Strategy.Operator = 0
	}

	bytes, _ := json.Marshal(sv.Publish)

	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ClusterId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId
	v.Data = bytes

	return v
}

func (s *strategyPublishVersionHandler) Decode(v *version_entry.Version) *strategy_entry.StrategyPublishVersion {
	sv := new(strategy_entry.StrategyPublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	val := make([]*strategy_entry.StrategyPublishConfigInfo, 0)
	_ = json.Unmarshal(v.Data, &val)
	sv.Publish = val
	return sv
}

func NewStrategyPublishVersionStore(db store.IDB, kind string) IStrategyPublishVersionStore {
	var h store.BaseKindHandler[strategy_entry.StrategyPublishVersion, version_entry.Version] = &strategyPublishVersionHandler{
		kind: kind,
	}
	return store.CreateBaseKindStore(h, db)
}
