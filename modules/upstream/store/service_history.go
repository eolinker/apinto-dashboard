package upstream_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceHistoryStore interface {
	store.BaseHistoryStore[entry.ServiceHistory]
}

type serviceHistoryHandler struct {
}

func (s *serviceHistoryHandler) Kind() string {
	return "service"
}

func (s *serviceHistoryHandler) Decode(r *entry.History) *entry.ServiceHistory {
	oldValue := new(entry.ServiceHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.ServiceHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.ServiceHistory{
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
	var historyHandler store.DecodeHistory[entry.ServiceHistory] = new(serviceHistoryHandler)
	return store.CreateHistory(historyHandler, db, entry.HistoryKindService)
}
