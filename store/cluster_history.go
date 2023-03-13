package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IClusterHistoryStore interface {
	BaseHistoryStore[entry.ClusterHistory]
}

type clusterHistoryHandler struct {
}

func (s *clusterHistoryHandler) Kind() string {
	return "cluster"
}

//func (s *clusterHistoryHandler) Encode(sr *entry.ClusterHistory) *entry.History {
//	oldValue, _ := json.Marshal(sr.OldValue)
//	newValue, _ := json.Marshal(sr.NewValue)
//	history := &entry.History{
//		Kind: s.Kind(),
//
//		NamespaceID: sr.NamespaceId,
//		TargetID:    sr.ClusterId,
//		OldValue:    string(oldValue),
//		NewValue:    string(newValue),
//		OptType:     sr.OptType,
//		OptTime:     sr.OptTime,
//		Operator:    sr.Operator,
//	}
//	return history
//}

func (s *clusterHistoryHandler) Decode(r *entry.History) *entry.ClusterHistory {
	oldValue := new(entry.Cluster)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.Cluster)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.ClusterHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ClusterId:   r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newClusterHistoryStore(db IDB) IClusterHistoryStore {
	var historyHandler DecodeHistory[entry.ClusterHistory] = new(clusterHistoryHandler)
	return CreateHistory(historyHandler, db, entry.HistoryKindCluster)
}
