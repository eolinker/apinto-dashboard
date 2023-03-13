package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IDiscoveryRuntimeStore interface {
	BaseRuntimeStore[entry.DiscoveryRuntime]
}

type discoveryRuntimeHandler struct {
}

func (s *discoveryRuntimeHandler) Kind() string {
	return "discovery"
}

func (s *discoveryRuntimeHandler) Encode(sr *entry.DiscoveryRuntime) *entry.Runtime {
	return &entry.Runtime{
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

func (s *discoveryRuntimeHandler) Decode(r *entry.Runtime) *entry.DiscoveryRuntime {
	return &entry.DiscoveryRuntime{
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

func newDiscoveryRuntimeStore(db IDB) IDiscoveryRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.DiscoveryRuntime, entry.Runtime] = new(discoveryRuntimeHandler)
	return CreateRuntime[entry.DiscoveryRuntime](runTimeHandler, db)

}
