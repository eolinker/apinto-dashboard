package store

import "github.com/eolinker/apinto-dashboard/entry"

type IClusterConfigRuntimeStore interface {
	BaseRuntimeStore[entry.ClusterConfigRuntime]
}

type clusterConfigRuntimeHandler struct {
}

func (s *clusterConfigRuntimeHandler) Kind() string {
	return "cluster_config"
}

func (s *clusterConfigRuntimeHandler) Encode(cr *entry.ClusterConfigRuntime) *entry.Runtime {
	return &entry.Runtime{
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

func (s *clusterConfigRuntimeHandler) Decode(r *entry.Runtime) *entry.ClusterConfigRuntime {
	return &entry.ClusterConfigRuntime{
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
	var runTimeHandler BaseKindHandler[entry.ClusterConfigRuntime, entry.Runtime] = new(clusterConfigRuntimeHandler)
	return createRuntime[entry.ClusterConfigRuntime](runTimeHandler, db)
}
