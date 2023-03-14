package api_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApiHistoryStore interface {
	store.BaseHistoryStore[api_entry.ApiHistory]
}

type ApiHistoryHandler struct {
}

func (s *ApiHistoryHandler) Kind() string {
	return "api"
}

func (s *ApiHistoryHandler) Decode(r *history_entry.History) *api_entry.ApiHistory {
	oldValue := new(api_entry.ApiHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(api_entry.ApiHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &api_entry.ApiHistory{
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

func newApiHistoryStore(db store.IDB) IApiHistoryStore {
	var historyHandler store.DecodeHistory[api_entry.ApiHistory] = new(ApiHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindAPI)
}
