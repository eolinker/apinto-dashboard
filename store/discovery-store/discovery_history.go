package discovery_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/discovery-entry"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDiscoveryHistoryStore interface {
	store.BaseHistoryStore[discovery_entry.DiscoveryHistory]
}

type discoveryHistoryHandler struct {
}

func (s *discoveryHistoryHandler) Kind() string {
	return "discovery"
}

func (s *discoveryHistoryHandler) Decode(r *history_entry.History) *discovery_entry.DiscoveryHistory {
	oldValue := new(discovery_entry.DiscoveryHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(discovery_entry.DiscoveryHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &discovery_entry.DiscoveryHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		DiscoveryId: r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newDiscoveryHistoryStore(db store.IDB) IDiscoveryHistoryStore {
	var historyHandler store.DecodeHistory[discovery_entry.DiscoveryHistory] = new(discoveryHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindDiscovery)
}
