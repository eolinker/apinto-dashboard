package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationRuntimeStore interface {
	BaseRuntimeStore[entry.ApplicationRuntime]
}

type applicationRuntimeHandler struct {
}

func (s *applicationRuntimeHandler) Kind() string {
	return "application"
}

func (s *applicationRuntimeHandler) Encode(sr *entry.ApplicationRuntime) *entry.Runtime {
	return &entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.ApplicationId,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionId,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		Disable:     sr.Disable,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *applicationRuntimeHandler) Decode(r *entry.Runtime) *entry.ApplicationRuntime {
	return &entry.ApplicationRuntime{
		Id:            r.Id,
		NamespaceId:   r.NamespaceID,
		ClusterId:     r.ClusterID,
		ApplicationId: r.TargetID,
		VersionId:     r.Version,
		IsOnline:      r.IsOnline,
		Disable:       r.Disable,
		Operator:      r.Operator,
		CreateTime:    r.CreateTime,
		UpdateTime:    r.UpdateTime,
	}
}

func newApplicationRuntimeStore(db IDB) IApplicationRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.ApplicationRuntime, entry.Runtime] = new(applicationRuntimeHandler)
	return createRuntime[entry.ApplicationRuntime](runTimeHandler, db)
}
