package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var _ IMonitorStore = (*monitorStore)(nil)

type IMonitorStore interface {
	IBaseStore[entry.MonitorPartition]
	GetByUUID(ctx context.Context, namespaceId int, uuid string) (*entry.MonitorPartition, error)
	GetByName(ctx context.Context, namespaceId int, name string) ([]*entry.MonitorPartition, error)
	GetList(ctx context.Context, namespaceId int) ([]*entry.MonitorPartition, error)
}

type monitorStore struct {
	*BaseStore[entry.MonitorPartition]
}

func (e *monitorStore) GetByUUID(ctx context.Context, namespaceId int, uuid string) (*entry.MonitorPartition, error) {
	return e.FirstQuery(ctx, "`namespace` = ? and `uuid` = ?", []interface{}{namespaceId, uuid}, "")
}

func (e *monitorStore) GetList(ctx context.Context, namespaceId int) ([]*entry.MonitorPartition, error) {
	return e.ListQuery(ctx, "`namespace` = ? ", []interface{}{namespaceId}, "create_time asc")
}

func (e *monitorStore) GetByName(ctx context.Context, namespaceId int, name string) ([]*entry.MonitorPartition, error) {
	return e.ListQuery(ctx, "`namespace` = ? and `name` = ? ", []interface{}{namespaceId, name}, "")
}

func newMonitorStore(db IDB) IMonitorStore {
	return &monitorStore{BaseStore: CreateStore[entry.MonitorPartition](db)}
}
