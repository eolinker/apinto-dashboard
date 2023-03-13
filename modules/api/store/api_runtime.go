package api_store

import (
	"github.com/eolinker/apinto-dashboard/entry"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IAPIRuntimeStore interface {
	store.BaseRuntimeStore[api_entry.APIRuntime]
}

type apiRuntimeHandler struct {
}

func (a *apiRuntimeHandler) Kind() string {
	return "api"
}

func (a *apiRuntimeHandler) Encode(ar *api_entry.APIRuntime) *entry.Runtime {
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

func (a *apiRuntimeHandler) Decode(r *entry.Runtime) *api_entry.APIRuntime {
	return &api_entry.APIRuntime{
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

func NewApiRuntimeStore(db store.IDB) IAPIRuntimeStore {
	var runTimeHandler store.BaseKindHandler[api_entry.APIRuntime, entry.Runtime] = new(apiRuntimeHandler)
	return store.CreateRuntime[api_entry.APIRuntime](runTimeHandler, db)
}
