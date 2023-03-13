package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApiHistoryStore interface {
	BaseHistoryStore[entry.ApiHistory]
}

type ApiHistoryHandler struct {
}

func (s *ApiHistoryHandler) Kind() string {
	return "api"
}

func (s *ApiHistoryHandler) Decode(r *entry.History) *entry.ApiHistory {
	oldValue := new(entry.ApiHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.ApiHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.ApiHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ApiId:       r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newApiHistoryStore(db IDB) IApiHistoryStore {
	var historyHandler DecodeHistory[entry.ApiHistory] = new(ApiHistoryHandler)
	return createHistory(historyHandler, db, entry.HistoryKindAPI)
}
