package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationHistoryStore interface {
	BaseHistoryStore[entry.ApplicationHistory]
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

func (s *applicationHistoryHandler) Decode(r *entry.History) *entry.ApplicationHistory {
	oldValue := new(entry.ApplicationHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.ApplicationHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.ApplicationHistory{
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

func newApplicationHistoryStore(db IDB) IApplicationHistoryStore {
	var historyHandler DecodeHistory[entry.ApplicationHistory] = new(applicationHistoryHandler)
	return createHistory[entry.ApplicationHistory](historyHandler, db, entry.HistoryKindApplication)
}
