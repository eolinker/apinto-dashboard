package store

import "github.com/eolinker/apinto-dashboard/entry"

type IAPIRuntimeStore interface {
	BaseRuntimeStore[entry.APIRuntime]
}

type apiRuntimeHandler struct {
}

func (a *apiRuntimeHandler) Kind() string {
	return "api"
}

func (a *apiRuntimeHandler) Encode(ar *entry.APIRuntime) *entry.Runtime {
	return &entry.Runtime{
		Id:          ar.Id,
		Kind:        a.Kind(),
		ClusterID:   ar.ClusterID,
		TargetID:    ar.ApiID,
		NamespaceID: ar.NamespaceId,
		Version:     ar.VersionID,
		IsOnline:    ar.IsOnline,
		Disable:     ar.Disable,
		Operator:    ar.Operator,
		CreateTime:  ar.CreateTime,
		UpdateTime:  ar.UpdateTime,
	}

}

func (a *apiRuntimeHandler) Decode(r *entry.Runtime) *entry.APIRuntime {
	return &entry.APIRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ApiID:       r.TargetID,
		ClusterID:   r.ClusterID,
		VersionID:   r.Version,
		IsOnline:    r.IsOnline,
		Disable:     r.Disable,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newApiRuntimeStore(db IDB) IAPIRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.APIRuntime, entry.Runtime] = new(apiRuntimeHandler)
	return createRuntime[entry.APIRuntime](runTimeHandler, db)
}
