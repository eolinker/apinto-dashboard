package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/publish-entry"
	"github.com/eolinker/apinto-dashboard/entry/strategy-entry"
)

type IStrategyPublishHistoryStore interface {
	BasePublishHistoryStore[strategy_entry.StrategyPublishHistory]
}

type strategyPublishHistoryHandler struct {
	kind string
}

func (s *strategyPublishHistoryHandler) Kind() string {
	return s.kind
}

func (s *strategyPublishHistoryHandler) Encode(sr *strategy_entry.StrategyPublishHistory) *publish_entry.PublishHistory {
	for _, publish := range sr.Publish {
		publish.Strategy.NamespaceId = 0
		publish.Strategy.ClusterId = 0
		publish.Strategy.Operator = 0
	}
	val, _ := json.Marshal(sr.Publish)
	history := &publish_entry.PublishHistory{
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

func (s *strategyPublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *strategy_entry.StrategyPublishHistory {
	val := make([]*strategy_entry.StrategyPublishConfigInfo, 0)
	_ = json.Unmarshal([]byte(r.Data), &val)
	history := &strategy_entry.StrategyPublishHistory{
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
	var historyHandler BaseKindHandler[strategy_entry.StrategyPublishHistory, publish_entry.PublishHistory] = &strategyPublishHistoryHandler{
		kind: kind,
	}
	return createPublishHistory[strategy_entry.StrategyPublishHistory](historyHandler, db)
}
