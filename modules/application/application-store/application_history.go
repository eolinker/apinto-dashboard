package application_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationHistoryStore interface {
	store.BaseHistoryStore[application_entry.ApplicationHistory]
}

type applicationHistoryHandler struct {
}

func (s *applicationHistoryHandler) Kind() string {
	return "application"
}

//func (s *applicationHistoryHandler) Encode(sr *entry.ApplicationHistory) *entry.History {
//	oldValue, _ := json.Marshal(sr.OldValue)
//	newValue, _ := json.Marshal(sr.NewValue)
//	history := &entry.History{
//		Kind:        s.Kind(),
//		ClusterID:   0,
//		NamespaceID: sr.NamespaceId,
//		TargetID:    sr.ApplicationId,
//		OldValue:    string(oldValue),
//		NewValue:    string(newValue),
//		OptType:     sr.OptType,
//		OptTime:     sr.OptTime,
//		Operator:    sr.Operator,
//	}
//	return history
//}

func (s *applicationHistoryHandler) Decode(r *history_entry.History) *application_entry.ApplicationHistory {
	oldValue := new(application_entry.ApplicationHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(application_entry.ApplicationHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &application_entry.ApplicationHistory{
		Id:            r.Id,
		NamespaceId:   r.NamespaceID,
		ApplicationId: r.TargetID,
		OptTime:       r.OptTime,
		OptType:       r.OptType,
		OldValue:      *oldValue,
		NewValue:      *newValue,
		Operator:      r.Operator,
	}
	return history
}

func newApplicationHistoryStore(db store.IDB) IApplicationHistoryStore {
	var historyHandler store.DecodeHistory[application_entry.ApplicationHistory] = new(applicationHistoryHandler)
	return store.CreateHistory[application_entry.ApplicationHistory](historyHandler, db, history_entry.HistoryKindApplication)
}
