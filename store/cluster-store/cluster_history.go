package cluster_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IClusterHistoryStore interface {
	store.BaseHistoryStore[cluster_entry.ClusterHistory]
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

func (s *clusterHistoryHandler) Decode(r *history_entry.History) *cluster_entry.ClusterHistory {
	oldValue := new(cluster_entry.Cluster)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(cluster_entry.Cluster)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &cluster_entry.ClusterHistory{
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

func newClusterHistoryStore(db store.IDB) IClusterHistoryStore {
	var historyHandler store.DecodeHistory[cluster_entry.ClusterHistory] = new(clusterHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindCluster)
}
