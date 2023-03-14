package store

import (
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
)

type IClusterConfigRuntimeStore interface {
	BaseRuntimeStore[cluster_entry.ClusterConfigRuntime]
}

type clusterConfigRuntimeHandler struct {
}

func (s *clusterConfigRuntimeHandler) Kind() string {
	return "cluster_config"
}

func (s *clusterConfigRuntimeHandler) Encode(cr *cluster_entry.ClusterConfigRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          cr.Id,
		Kind:        s.Kind(),
		ClusterID:   cr.ClusterId,
		TargetID:    cr.ConfigID,
		NamespaceID: cr.NamespaceId,
		Version:     0,
		IsOnline:    cr.IsOnline,
		Disable:     false,
		Operator:    cr.Operator,
		CreateTime:  cr.CreateTime,
		UpdateTime:  cr.UpdateTime,
	}

}

func (s *clusterConfigRuntimeHandler) Decode(r *runtime_entry.Runtime) *cluster_entry.ClusterConfigRuntime {
	return &cluster_entry.ClusterConfigRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ConfigID:    r.TargetID,
		ClusterId:   r.ClusterID,
		IsOnline:    r.IsOnline,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newClusterConfigRuntimeStore(db IDB) IClusterConfigRuntimeStore {
	var runTimeHandler BaseKindHandler[cluster_entry.ClusterConfigRuntime, runtime_entry.Runtime] = new(clusterConfigRuntimeHandler)
	return CreateRuntime[cluster_entry.ClusterConfigRuntime](runTimeHandler, db)
}
