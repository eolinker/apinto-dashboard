package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/base/publish-entry"
)

type BasePublishHistoryStore[T any] interface {
	IBaseStore[T]
	GetByVersionName(ctx context.Context, versionName string, targetId int) (*T, error)
	GetByClusterPage(ctx context.Context, pageNum, pageSize, clusterId int) ([]*T, int, error)
}

type BasePublishHistory[T any] struct {
	*BaseKindStore[T, publish_entry.PublishHistory]
}

func CreatePublishHistory[T any](handler BaseKindHandler[T, publish_entry.PublishHistory], db IDB) *BasePublishHistory[T] {
	return &BasePublishHistory[T]{
		BaseKindStore: CreateBaseKindStore[T, publish_entry.PublishHistory](handler, db),
	}
}

func (b *BasePublishHistory[T]) GetByVersionName(ctx context.Context, versionName string, targetId int) (*T, error) {
	return b.First(ctx, map[string]interface{}{"`kind`": b.BaseKindHandler.Kind(), "`target`": targetId, "version_name": versionName})
}

func (b *BasePublishHistory[T]) GetByClusterPage(ctx context.Context, pageNum, pageSize, clusterId int) ([]*T, int, error) {
	return b.ListPage(ctx, "`kind` = ? and `cluster` = ?", pageNum, pageSize, []interface{}{b.BaseKindHandler.Kind(), clusterId}, "")
}
