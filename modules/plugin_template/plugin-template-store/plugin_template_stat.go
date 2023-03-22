package plugin_template_store

import (
	stat_entry "github.com/eolinker/apinto-dashboard/modules/base/stat-entry"
	plugin_template_entry "github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginTemplateStatStore interface {
	store.IBaseStore[plugin_template_entry.PluginTemplateStat]
}

type pluginTemplateHandler struct {
}

func (p *pluginTemplateHandler) Kind() string {
	return "plugin_template"
}

func (p *pluginTemplateHandler) Encode(as *plugin_template_entry.PluginTemplateStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = as.PluginTemplateId
	stat.Kind = p.Kind()
	stat.Version = as.VersionID

	return stat
}

func (p *pluginTemplateHandler) Decode(stat *stat_entry.Stat) *plugin_template_entry.PluginTemplateStat {
	ds := new(plugin_template_entry.PluginTemplateStat)

	ds.PluginTemplateId = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newPluginTemplateStatStore(db store.IDB) IPluginTemplateStatStore {
	var h store.BaseKindHandler[plugin_template_entry.PluginTemplateStat, stat_entry.Stat] = new(pluginTemplateHandler)
	return store.CreateBaseKindStore(h, db)
}
