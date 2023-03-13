package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IDiscoveryHistoryStore interface {
	BaseHistoryStore[entry.DiscoveryHistory]
}

type discoveryHistoryHandler struct {
}

func (s *discoveryHistoryHandler) Kind() string {
	return "discovery"
}

func (s *discoveryHistoryHandler) Decode(r *entry.History) *entry.DiscoveryHistory {
	oldValue := new(entry.DiscoveryHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.DiscoveryHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.DiscoveryHistory{
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

func newDiscoveryHistoryStore(db IDB) IDiscoveryHistoryStore {
	var historyHandler DecodeHistory[entry.DiscoveryHistory] = new(discoveryHistoryHandler)
	return CreateHistory(historyHandler, db, entry.HistoryKindDiscovery)
}
