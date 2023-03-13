package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IVariableRuntimeStore interface {
	BaseRuntimeStore[entry.VariableRuntime]
}

type variableRuntimeHandler struct {
}

func (s *variableRuntimeHandler) Kind() string {
	return "variable"
}

func (s *variableRuntimeHandler) Encode(sr *entry.VariableRuntime) *entry.Runtime {
	return &entry.Runtime{
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

func (s *variableRuntimeHandler) Decode(r *entry.Runtime) *entry.VariableRuntime {
	return &entry.VariableRuntime{
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

func newVariableRuntimeStore(db IDB) IVariableRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.VariableRuntime, entry.Runtime] = new(variableRuntimeHandler)
	return CreateRuntime[entry.VariableRuntime](runTimeHandler, db)
}
