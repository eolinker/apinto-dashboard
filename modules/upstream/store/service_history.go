package upstream_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"github.com/eolinker/apinto-dashboard/entry/upstream-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceHistoryStore interface {
	store.BaseHistoryStore[upstream_entry.ServiceHistory]
}

type serviceHistoryHandler struct {
}

func (s *serviceHistoryHandler) Kind() string {
	return "service"
}

func (s *serviceHistoryHandler) Decode(r *history_entry.History) *upstream_entry.ServiceHistory {
	oldValue := new(upstream_entry.ServiceHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(upstream_entry.ServiceHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &upstream_entry.ServiceHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ServiceId:   r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newServiceHistoryStore(db store.IDB) IServiceHistoryStore {
	var historyHandler store.DecodeHistory[upstream_entry.ServiceHistory] = new(serviceHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindService)
}
