package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
)

type BaseRuntimeStore[T any] interface {
	IBaseStore[T]
	GetByTarget(ctx context.Context, id int) ([]*T, error)
	GetByCluster(ctx context.Context, clusterId int) ([]*T, error)
	GetForCluster(ctx context.Context, id int, clusterId int) (*T, error)
	OnlineCount(ctx context.Context, id int) (int64, error)
}

type BaseRuntime[T any] struct {
	*BaseKindStore[T, runtime_entry.Runtime]
}

func (b *BaseRuntime[T]) GetByTarget(ctx context.Context, target int) ([]*T, error) {
	return b.BaseKindStore.List(ctx, map[string]interface{}{
		"kind":   b.BaseKindHandler.Kind(),
		"target": target,
	})
}

func (b *BaseRuntime[T]) GetByCluster(ctx context.Context, clusterId int) ([]*T, error) {
	return b.BaseKindStore.List(ctx, map[string]interface{}{
		"cluster": clusterId,
		"kind":    b.BaseKindHandler.Kind(),
	})
}

func (b *BaseRuntime[T]) GetForCluster(ctx context.Context, target int, clusterId int) (*T, error) {
	return b.BaseKindStore.First(ctx, map[string]interface{}{
		"kind":    b.BaseKindHandler.Kind(),
		"target":  target,
		"cluster": clusterId,
	})
}

func (b *BaseRuntime[T]) OnlineCount(ctx context.Context, target int) (int64, error) {
	count := int64(0)
	err := b.DB(ctx).Table("`runtime`").Where("`kind` = ? and `target` = ? and `is_online` = 1", b.BaseKindHandler.Kind(), target).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CreateRuntime[T any](handler BaseKindHandler[T, runtime_entry.Runtime], db IDB) *BaseRuntime[T] {
	return &BaseRuntime[T]{
		BaseKindStore: CreateBaseKindStore[T, runtime_entry.Runtime](handler, db),
	}
}
