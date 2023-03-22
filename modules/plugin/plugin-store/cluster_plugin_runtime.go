package plugin_store

import (
	"github.com/eolinker/apinto-dashboard/modules/base/runtime-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IClusterPluginRuntimeStore interface {
	store.BaseRuntimeStore[plugin_entry.ClusterPluginRuntime]
}

type clusterPluginRuntimeHandler struct {
}

func (s *clusterPluginRuntimeHandler) Kind() string {
	return "cluster_plugin"
}

func (s *clusterPluginRuntimeHandler) Encode(sr *plugin_entry.ClusterPluginRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.ClusterId,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionId,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *clusterPluginRuntimeHandler) Decode(r *runtime_entry.Runtime) *plugin_entry.ClusterPluginRuntime {
	return &plugin_entry.ClusterPluginRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ClusterId:   r.ClusterID,
		VersionId:   r.Version,
		IsOnline:    r.IsOnline,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newClusterPluginRuntimeStore(db store.IDB) IClusterPluginRuntimeStore {
	var runTimeHandler store.BaseKindHandler[plugin_entry.ClusterPluginRuntime, runtime_entry.Runtime] = new(clusterPluginRuntimeHandler)
	return store.CreateRuntime[plugin_entry.ClusterPluginRuntime](runTimeHandler, db)
}
