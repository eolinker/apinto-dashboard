package plugin_template_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginTemplateHistoryStore interface {
	store.BaseHistoryStore[plugin_template_entry.PluginTemplateHistory]
}

type PluginTemplateHistoryHandler struct {
}

func (s *PluginTemplateHistoryHandler) Kind() string {
	return "PluginTemplate"
}

func (s *PluginTemplateHistoryHandler) Decode(r *history_entry.History) *plugin_template_entry.PluginTemplateHistory {
	oldValue := new(plugin_template_entry.PluginTemplateHistoryInfo)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(plugin_template_entry.PluginTemplateHistoryInfo)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &plugin_template_entry.PluginTemplateHistory{
		Id:               r.Id,
		NamespaceId:      r.NamespaceID,
		PluginTemplateID: r.TargetID,
		OptTime:          r.OptTime,
		OptType:          r.OptType,
		OldValue:         *oldValue,
		NewValue:         *newValue,
		Operator:         r.Operator,
	}
	return history
}

func newPluginTemplateHistoryStore(db store.IDB) IPluginTemplateHistoryStore {
	var historyHandler store.DecodeHistory[plugin_template_entry.PluginTemplateHistory] = new(PluginTemplateHistoryHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindPluginTemplate)
}
