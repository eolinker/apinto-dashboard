package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IStrategyHistoryStore interface {
	BaseHistoryStore[entry.StrategyHistory]
}

type strategyHistoryHandler struct {
}

func (s *strategyHistoryHandler) Kind() string {
	return "strategy"
}

func (s *strategyHistoryHandler) Decode(r *entry.History) *entry.StrategyHistory {
	oldValue := new(entry.StrategyHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.StrategyHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.StrategyHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		StrategyId:  r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newStrategyHistoryStore(db IDB) IStrategyHistoryStore {
	var historyHandler DecodeHistory[entry.StrategyHistory] = new(strategyHistoryHandler)
	return CreateHistory(historyHandler, db, entry.HistoryKindStrategy)
}
