package upstream_store

import (
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
	"github.com/eolinker/apinto-dashboard/entry/upstream-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceRuntimeStore interface {
	store.BaseRuntimeStore[upstream_entry.ServiceRuntime]
}

type serviceRuntimeHandler struct {
}

func (s *serviceRuntimeHandler) Kind() string {
	return "service"
}

func (s *serviceRuntimeHandler) Encode(sr *upstream_entry.ServiceRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.ServiceId,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionId,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *serviceRuntimeHandler) Decode(r *runtime_entry.Runtime) *upstream_entry.ServiceRuntime {
	return &upstream_entry.ServiceRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ServiceId:   r.TargetID,
		ClusterId:   r.ClusterID,
		VersionId:   r.Version,
		IsOnline:    r.IsOnline,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newServiceRuntimeStore(db store.IDB) IServiceRuntimeStore {
	var runTimeHandler store.BaseKindHandler[upstream_entry.ServiceRuntime, runtime_entry.Runtime] = new(serviceRuntimeHandler)
	return store.CreateRuntime[upstream_entry.ServiceRuntime](runTimeHandler, db)

}
