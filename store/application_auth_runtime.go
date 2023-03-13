package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthRuntimeStore interface {
	BaseRuntimeStore[entry.ApplicationAuthRuntime]
}

type applicationAuthRuntimeHandler struct {
}

func (s *applicationAuthRuntimeHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthRuntimeHandler) Encode(sr *entry.ApplicationAuthRuntime) *entry.Runtime {
	return &entry.Runtime{
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

func (s *applicationAuthRuntimeHandler) Decode(r *entry.Runtime) *entry.ApplicationAuthRuntime {
	return &entry.ApplicationAuthRuntime{
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

func newApplicationAuthRuntimeStore(db IDB) IApplicationAuthRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.ApplicationAuthRuntime, entry.Runtime] = new(applicationAuthRuntimeHandler)
	return createRuntime[entry.ApplicationAuthRuntime](runTimeHandler, db)
}
