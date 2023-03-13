package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthHistoryStore interface {
	BaseHistoryStore[entry.ApplicationAuthHistory]
}

type applicationAuthHistoryHandler struct {
}

func (s *applicationAuthHistoryHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthHistoryHandler) Decode(r *entry.History) *entry.ApplicationAuthHistory {
	oldValue := new(entry.ApplicationAuthHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.ApplicationAuthHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.ApplicationAuthHistory{
		Id:                r.Id,
		NamespaceId:       r.NamespaceID,
		ApplicationAuthId: r.TargetID,
		OptTime:           r.OptTime,
		OptType:           r.OptType,
		OldValue:          *oldValue,
		NewValue:          *newValue,
		Operator:          r.Operator,
	}
	return history
}

func newApplicationAuthHistoryStore(db IDB) IApplicationAuthHistoryStore {
	var historyHandler DecodeHistory[entry.ApplicationAuthHistory] = new(applicationAuthHistoryHandler)
	return CreateHistory(historyHandler, db, entry.HistoryKindApplicationAuth)
}
