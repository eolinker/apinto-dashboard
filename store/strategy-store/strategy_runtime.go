package strategy_store

import (
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
	"github.com/eolinker/apinto-dashboard/entry/strategy-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

//流量策略

type IStrategyRuntimeStore interface {
	store.BaseRuntimeStore[strategy_entry.StrategyRuntime]
}

type strategyRuntimeHandler struct {
	kind string
}

func (s *strategyRuntimeHandler) Kind() string {
	return s.kind
}

func (s *strategyRuntimeHandler) Encode(sr *strategy_entry.StrategyRuntime) *runtime_entry.Runtime {
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

func (s *strategyRuntimeHandler) Decode(r *runtime_entry.Runtime) *strategy_entry.StrategyRuntime {
	return &strategy_entry.StrategyRuntime{
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

func NewStrategyRuntimeStore(db store.IDB, kind string) IStrategyRuntimeStore {
	var runTimeHandler store.BaseKindHandler[strategy_entry.StrategyRuntime, runtime_entry.Runtime] = &strategyRuntimeHandler{
		kind: kind,
	}
	return store.CreateRuntime[strategy_entry.StrategyRuntime](runTimeHandler, db)
}
