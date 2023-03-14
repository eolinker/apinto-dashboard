package variable_store

import (
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IVariableRuntimeStore interface {
	store.BaseRuntimeStore[variable_entry.VariableRuntime]
}

type variableRuntimeHandler struct {
}

func (s *variableRuntimeHandler) Kind() string {
	return "variable"
}

func (s *variableRuntimeHandler) Encode(sr *variable_entry.VariableRuntime) *runtime_entry.Runtime {
	return &runtime_entry.Runtime{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterID:   sr.ClusterId,
		TargetID:    sr.ClusterId,
		NamespaceID: sr.NamespaceId,
		Version:     sr.VersionId,
		IsOnline:    sr.IsOnline,
		Operator:    sr.Operator,
		CreateTime:  sr.CreateTime,
		UpdateTime:  sr.UpdateTime,
	}

}

func (s *variableRuntimeHandler) Decode(r *runtime_entry.Runtime) *variable_entry.VariableRuntime {
	return &variable_entry.VariableRuntime{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		ClusterId:   r.ClusterID,
		VersionId:   r.Version,
		IsOnline:    r.IsOnline,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}

func newVariableRuntimeStore(db store.IDB) IVariableRuntimeStore {
	var runTimeHandler store.BaseKindHandler[variable_entry.VariableRuntime, runtime_entry.Runtime] = new(variableRuntimeHandler)
	return store.CreateRuntime[variable_entry.VariableRuntime](runTimeHandler, db)
}
