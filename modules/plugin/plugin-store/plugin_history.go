package plugin_store

import (
	"encoding/json"
	history_entry "github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginHistoryStore interface {
	store.BaseHistoryStore[plugin_entry.PluginHistory]
}

type PluginHistoryHandler struct {
}

func (s *PluginHistoryHandler) Kind() string {
	return "plugin"
}

func (s *PluginHistoryHandler) Decode(r *history_entry.History) *plugin_entry.PluginHistory {
	oldValue := new(plugin_entry.PluginHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(plugin_entry.PluginHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &plugin_entry.PluginHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		PluginID:    r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newPluginHistoryStore(db store.IDB) IPluginHistoryStore {
	var historyHandler store.DecodeHistory[plugin_entry.PluginHistory] = new(PluginHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindPlugin)
}
