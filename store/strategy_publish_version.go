package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IStrategyPublishVersionStore interface {
	IBaseStore[entry.StrategyPublishVersion]
}

type strategyPublishVersionHandler struct {
	kind string
}

func (s *strategyPublishVersionHandler) Kind() string {
	return s.kind
}

func (s *strategyPublishVersionHandler) Encode(sv *entry.StrategyPublishVersion) *entry.Version {

	for _, publish := range sv.Publish {
		publish.Strategy.NamespaceId = 0
		publish.Strategy.ClusterId = 0
		publish.Strategy.Operator = 0
	}

	bytes, _ := json.Marshal(sv.Publish)

	v := new(entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ClusterId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId
	v.Data = bytes

	return v
}

func (s *strategyPublishVersionHandler) Decode(v *entry.Version) *entry.StrategyPublishVersion {
	sv := new(entry.StrategyPublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	val := make([]*entry.StrategyPublishConfigInfo, 0)
	_ = json.Unmarshal(v.Data, &val)
	sv.Publish = val
	return sv
}

func NewStrategyPublishVersionStore(db IDB, kind string) IStrategyPublishVersionStore {
	var h BaseKindHandler[entry.StrategyPublishVersion, entry.Version] = &strategyPublishVersionHandler{
		kind: kind,
	}
	return CreateBaseKindStore(h, db)
}
