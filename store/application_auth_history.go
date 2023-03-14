package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
)

type IApplicationAuthHistoryStore interface {
	BaseHistoryStore[application_entry.ApplicationAuthHistory]
}

type applicationAuthHistoryHandler struct {
}

func (s *applicationAuthHistoryHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthHistoryHandler) Decode(r *history_entry.History) *application_entry.ApplicationAuthHistory {
	oldValue := new(application_entry.ApplicationAuthHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(application_entry.ApplicationAuthHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &application_entry.ApplicationAuthHistory{
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
	var historyHandler DecodeHistory[application_entry.ApplicationAuthHistory] = new(applicationAuthHistoryHandler)
	return CreateHistory(historyHandler, db, history_entry.HistoryKindApplicationAuth)
}
