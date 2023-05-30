package plugin_template_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/publish-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginTemplatePublishHistoryStore interface {
	store.BasePublishHistoryStore[plugin_template_entry.PluginTemplatePublishHistory]
}

type pluginTemplatePublishHistoryHandler struct {
}

func (s *pluginTemplatePublishHistoryHandler) Kind() string {
	return "plugin_template"
}

func (s *pluginTemplatePublishHistoryHandler) Encode(sr *plugin_template_entry.PluginTemplatePublishHistory) *publish_entry.PublishHistory {
	val, _ := json.Marshal(sr.PluginTemplateVersionConfig)
	history := &publish_entry.PublishHistory{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterId:   sr.ClusterId,
		NamespaceId: sr.NamespaceId,
		Target:      sr.Target,
		VersionId:   sr.VersionId,
		Data:        string(val),
		Desc:        sr.Desc,
		OptType:     sr.OptType,
		OptTime:     sr.OptTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *pluginTemplatePublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *plugin_template_entry.PluginTemplatePublishHistory {
	val := new(plugin_template_entry.PluginTemplateVersionConfig)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &plugin_template_entry.PluginTemplatePublishHistory{
		Id:                          r.Id,
		VersionName:                 r.VersionName,
		ClusterId:                   r.ClusterId,
		NamespaceId:                 r.NamespaceId,
		Desc:                        r.Desc,
		VersionId:                   r.VersionId,
		Target:                      r.Target,
		PluginTemplateVersionConfig: *val,
		OptType:                     r.OptType,
		Operator:                    r.Operator,
		OptTime:                     r.OptTime,
	}
	return history
}

func newPluginTemplatePublishHistoryStore(db store.IDB) IPluginTemplatePublishHistoryStore {
	var historyHandler store.BaseKindHandler[plugin_template_entry.PluginTemplatePublishHistory, publish_entry.PublishHistory] = new(pluginTemplatePublishHistoryHandler)
	return store.CreatePublishHistory[plugin_template_entry.PluginTemplatePublishHistory](historyHandler, db)
}
