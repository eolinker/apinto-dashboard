package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IServiceRuntimeStore interface {
	BaseRuntimeStore[entry.ServiceRuntime]
}

type serviceRuntimeHandler struct {
}

func (s *serviceRuntimeHandler) Kind() string {
	return "service"
}

func (s *serviceRuntimeHandler) Encode(sr *entry.ServiceRuntime) *entry.Runtime {
	return &entry.Runtime{
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

func (s *serviceRuntimeHandler) Decode(r *entry.Runtime) *entry.ServiceRuntime {
	return &entry.ServiceRuntime{
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

func newServiceRuntimeStore(db IDB) IServiceRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.ServiceRuntime, entry.Runtime] = new(serviceRuntimeHandler)
	return createRuntime[entry.ServiceRuntime](runTimeHandler, db)

}
