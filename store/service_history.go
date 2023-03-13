package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IServiceHistoryStore interface {
	BaseHistoryStore[entry.ServiceHistory]
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

func newServiceHistoryStore(db IDB) IServiceHistoryStore {
	var historyHandler DecodeHistory[entry.ServiceHistory] = new(serviceHistoryHandler)
	return createHistory(historyHandler, db, entry.HistoryKindService)
}
