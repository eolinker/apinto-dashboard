package application_store

import (
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationAuthRuntimeStore interface {
	store.BaseRuntimeStore[application_entry.ApplicationAuthRuntime]
}

type applicationAuthRuntimeHandler struct {
}

func (s *applicationAuthRuntimeHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthRuntimeHandler) Encode(sr *application_entry.ApplicationAuthRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.ApplicationAuthId,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionId,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		Disable:     sr.Disable,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *applicationAuthRuntimeHandler) Decode(r *runtime_entry.Runtime) *application_entry.ApplicationAuthRuntime {
	return &application_entry.ApplicationAuthRuntime{
		Id:                r.Id,
		NamespaceId:       r.NamespaceID,
		ClusterId:         r.ClusterID,
		ApplicationAuthId: r.TargetID,
		VersionId:         r.Version,
		IsOnline:          r.IsOnline,
		Disable:           r.Disable,
		Operator:          r.Operator,
		CreateTime:        r.CreateTime,
		UpdateTime:        r.UpdateTime,
	}
}

func newApplicationAuthRuntimeStore(db store.IDB) IApplicationAuthRuntimeStore {
	var runTimeHandler store.BaseKindHandler[application_entry.ApplicationAuthRuntime, runtime_entry.Runtime] = new(applicationAuthRuntimeHandler)
	return store.CreateRuntime[application_entry.ApplicationAuthRuntime](runTimeHandler, db)
}
