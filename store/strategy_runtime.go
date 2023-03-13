package store

import "github.com/eolinker/apinto-dashboard/entry"

//流量策略

type IStrategyRuntimeStore interface {
	BaseRuntimeStore[entry.StrategyRuntime]
}

type strategyRuntimeHandler struct {
	kind string
}

func (s *strategyRuntimeHandler) Kind() string {
	return s.kind
}

func (s *strategyRuntimeHandler) Encode(sr *entry.StrategyRuntime) *entry.Runtime {
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

func (s *strategyRuntimeHandler) Decode(r *entry.Runtime) *entry.StrategyRuntime {
	return &entry.StrategyRuntime{
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

func NewStrategyRuntimeStore(db IDB, kind string) IStrategyRuntimeStore {
	var runTimeHandler BaseKindHandler[entry.StrategyRuntime, entry.Runtime] = &strategyRuntimeHandler{
		kind: kind,
	}
	return CreateRuntime[entry.StrategyRuntime](runTimeHandler, db)
}
