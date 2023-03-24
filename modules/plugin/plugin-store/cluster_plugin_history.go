package plugin_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IClusterPluginHistoryStore interface {
	store.BaseHistoryStore[plugin_entry.ClusterPluginHistory]
}

type clusterPluginHistoryHandler struct {
}

func (c *clusterPluginHistoryHandler) Kind() string {
	return "cluster_plugin"
}

func (c *clusterPluginHistoryHandler) Decode(r *history_entry.History) *plugin_entry.ClusterPluginHistory {
	oldValue := new(plugin_entry.ClusterPluginHistoryValue)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(plugin_entry.ClusterPluginHistoryValue)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &plugin_entry.ClusterPluginHistory{
		Id:              r.Id,
		NamespaceId:     r.NamespaceID,
		ClusterPluginID: r.TargetID,
		OldValue:        *oldValue,
		NewValue:        *newValue,
		OptType:         r.OptType,
		Operator:        r.Operator,
		OptTime:         r.OptTime,
	}
	return history
}

func newClusterPluginHistoryStore(db store.IDB) IClusterPluginHistoryStore {
	var historyHandler store.DecodeHistory[plugin_entry.ClusterPluginHistory] = new(clusterPluginHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindClusterPlugin)
}
