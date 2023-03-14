package store

import (
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
)

type IApplicationRuntimeStore interface {
	BaseRuntimeStore[application_entry.ApplicationRuntime]
}

type applicationRuntimeHandler struct {
}

func (s *applicationRuntimeHandler) Kind() string {
	return "application"
}

func (s *applicationRuntimeHandler) Encode(sr *application_entry.ApplicationRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
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

func (s *applicationRuntimeHandler) Decode(r *runtime_entry.Runtime) *application_entry.ApplicationRuntime {
	return &application_entry.ApplicationRuntime{
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
	var runTimeHandler BaseKindHandler[application_entry.ApplicationRuntime, runtime_entry.Runtime] = new(applicationRuntimeHandler)
	return CreateRuntime[application_entry.ApplicationRuntime](runTimeHandler, db)
}
