package plugin_template_store

import (
	"github.com/eolinker/apinto-dashboard/modules/base/runtime-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginTemplateRuntimeStore interface {
	store.BaseRuntimeStore[plugin_template_entry.PluginTemplateRuntime]
}

type PluginTemplateRuntimeHandler struct {
}

func (a *PluginTemplateRuntimeHandler) Kind() string {
	return "plugin_template"
}

func (a *PluginTemplateRuntimeHandler) Encode(ar *plugin_template_entry.PluginTemplateRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          ar.Id,
		Kind:        a.Kind(),
		ClusterID:   ar.ClusterID,
		TargetID:    ar.PluginTemplateID,
		NamespaceID: ar.NamespaceId,
		Version:     ar.VersionID,
		IsOnline:    ar.IsOnline,
		Disable:     ar.Disable,
		Operator:    ar.Operator,
		CreateTime:  ar.CreateTime,
		UpdateTime:  ar.UpdateTime,
	}

}

func (a *PluginTemplateRuntimeHandler) Decode(r *runtime_entry.Runtime) *plugin_template_entry.PluginTemplateRuntime {
	return &plugin_template_entry.PluginTemplateRuntime{
		Id:               r.Id,
		NamespaceId:      r.NamespaceID,
		PluginTemplateID: r.TargetID,
		ClusterID:        r.ClusterID,
		VersionID:        r.Version,
		IsOnline:         r.IsOnline,
		Disable:          r.Disable,
		Operator:         r.Operator,
		CreateTime:       r.CreateTime,
		UpdateTime:       r.UpdateTime,
	}
}

func newPluginTemplateRuntimeStore(db store.IDB) IPluginTemplateRuntimeStore {
	var runTimeHandler store.BaseKindHandler[plugin_template_entry.PluginTemplateRuntime, runtime_entry.Runtime] = new(PluginTemplateRuntimeHandler)
	return store.CreateRuntime[plugin_template_entry.PluginTemplateRuntime](runTimeHandler, db)
}
