package plugin_template_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/version-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginTemplateVersionStore interface {
	store.IBaseStore[plugin_template_entry.PluginTemplateVersion]
}

type PluginTemplateVersionStore struct {
	*store.BaseKindStore[plugin_template_entry.PluginTemplateVersion, version_entry.Version]
}

type pluginTemplateVersionKindHandler struct {
}

func (p *pluginTemplateVersionKindHandler) Kind() string {
	return "plugin_template"
}

func (p *pluginTemplateVersionKindHandler) Encode(av *plugin_template_entry.PluginTemplateVersion) *version_entry.Version {
	data, _ := json.Marshal(av.PluginTemplateVersionConfig)
	v := &version_entry.Version{
		Id:          av.Id,
		Target:      av.PluginTemplateId,
		NamespaceID: av.NamespaceID,
		Kind:        p.Kind(),
		Data:        data,
		Operator:    av.Operator,
		CreateTime:  av.CreateTime,
	}

	return v
}

func (p *pluginTemplateVersionKindHandler) Decode(v *version_entry.Version) *plugin_template_entry.PluginTemplateVersion {
	av := &plugin_template_entry.PluginTemplateVersion{
		Id:                          v.Id,
		PluginTemplateId:            v.Target,
		NamespaceID:                 v.NamespaceID,
		PluginTemplateVersionConfig: plugin_template_entry.PluginTemplateVersionConfig{},
		Operator:                    v.Operator,
		CreateTime:                  v.CreateTime,
	}
	_ = json.Unmarshal(v.Data, &av.PluginTemplateVersionConfig)

	return av
}

func newPluginTemplateVersionStore(db store.IDB) IPluginTemplateVersionStore {
	var h store.BaseKindHandler[plugin_template_entry.PluginTemplateVersion, version_entry.Version] = new(pluginTemplateVersionKindHandler)
	return &PluginTemplateVersionStore{BaseKindStore: store.CreateBaseKindStore(h, db)}
}
