package discovery_store

import (
	"github.com/eolinker/apinto-dashboard/modules/base/runtime-entry"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discovery-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDiscoveryRuntimeStore interface {
	store.BaseRuntimeStore[discovery_entry.DiscoveryRuntime]
}

type discoveryRuntimeHandler struct {
}

func (s *discoveryRuntimeHandler) Kind() string {
	return "discovery"
}

func (s *discoveryRuntimeHandler) Encode(sr *discovery_entry.DiscoveryRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.DiscoveryID,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionID,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *discoveryRuntimeHandler) Decode(r *runtime_entry.Runtime) *discovery_entry.DiscoveryRuntime {
	return &discovery_entry.DiscoveryRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		DiscoveryID: r.TargetID,
		ClusterId:   r.ClusterID,
		VersionID:   r.Version,
		IsOnline:    r.IsOnline,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newDiscoveryRuntimeStore(db store.IDB) IDiscoveryRuntimeStore {
	var runTimeHandler store.BaseKindHandler[discovery_entry.DiscoveryRuntime, runtime_entry.Runtime] = new(discoveryRuntimeHandler)
	return store.CreateRuntime[discovery_entry.DiscoveryRuntime](runTimeHandler, db)

}
