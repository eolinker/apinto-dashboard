package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IStrategyPublishHistoryStore interface {
	BasePublishHistoryStore[entry.StrategyPublishHistory]
}

type strategyPublishHistoryHandler struct {
	kind string
}

func (s *strategyPublishHistoryHandler) Kind() string {
	return s.kind
}

func (s *strategyPublishHistoryHandler) Encode(sr *entry.StrategyPublishHistory) *entry.PublishHistory {
	for _, publish := range sr.Publish {
		publish.Strategy.NamespaceId = 0
		publish.Strategy.ClusterId = 0
		publish.Strategy.Operator = 0
	}
	val, _ := json.Marshal(sr.Publish)
	history := &entry.PublishHistory{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterId:   sr.ClusterId,
		NamespaceId: sr.NamespaceId,
		Target:      sr.ClusterId,
		VersionId:   sr.VersionId,
		Data:        string(val),
		Desc:        sr.Desc,
		OptType:     sr.OptType,
		OptTime:     sr.CreateTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *strategyPublishHistoryHandler) Decode(r *entry.PublishHistory) *entry.StrategyPublishHistory {
	val := make([]*entry.StrategyPublishConfigInfo, 0)
	_ = json.Unmarshal([]byte(r.Data), &val)
	history := &entry.StrategyPublishHistory{
		Id:          r.Id,
		VersionName: r.VersionName,
		Desc:        r.Desc,
		NamespaceId: r.NamespaceId,
		ClusterId:   r.ClusterId,
		VersionId:   r.VersionId,
		CreateTime:  r.OptTime,
		OptType:     r.OptType,
		Publish:     val,
		Operator:    r.Operator,
	}
	return history
}

func NewStrategyPublishHistoryStore(db IDB, kind string) IStrategyPublishHistoryStore {
	var historyHandler BaseKindHandler[entry.StrategyPublishHistory, entry.PublishHistory] = &strategyPublishHistoryHandler{
		kind: kind,
	}
	return createPublishHistory[entry.StrategyPublishHistory](historyHandler, db)
}
