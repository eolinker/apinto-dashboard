package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var _ IApplicationStore = (*applicationStore)(nil)

type IApplicationStore interface {
	IBaseStore[entry.Application]
	GetByIdStr(ctx context.Context, namespaceId int, idStr string) (*entry.Application, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*entry.Application, error)
	GetListPage(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*entry.Application, int, error)
	GetList(ctx context.Context, namespaceId int, uuids ...string) ([]*entry.Application, error)
	GetListByName(ctx context.Context, namespaceId int, name string) ([]*entry.Application, error)
}

type applicationStore struct {
	*BaseStore[entry.Application]
}

func (a *applicationStore) GetByIdStr(ctx context.Context, namespaceId int, idStr string) (*entry.Application, error) {
	return a.FirstQuery(ctx, "`namespace` = ? and `id_str` = ?", []interface{}{namespaceId, idStr}, "")
}

func (a *applicationStore) GetByName(ctx context.Context, namespaceId int, name string) (*entry.Application, error) {
	return a.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (a *applicationStore) GetListPage(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*entry.Application, int, error) {
	if queryName == "" {
		return a.ListPage(ctx, "`namespace` = ?", pageNum, pageSize, []interface{}{namespaceId}, "update_time desc")
	}
	return a.ListPage(ctx, "`namespace` = ? and `name` like ?", pageNum, pageSize, []interface{}{namespaceId, "%" + queryName + "%"}, "update_time desc")
}

func (a *applicationStore) GetList(ctx context.Context, namespaceId int, uuids ...string) ([]*entry.Application, error) {
	if len(uuids) > 0 {
		return a.ListQuery(ctx, "`namespace` = ? and `id_str` in (?)", []interface{}{namespaceId, uuids}, "create_time asc")
	}
	return a.ListQuery(ctx, "`namespace` = ?", []interface{}{namespaceId}, "create_time asc")
}

func (a *applicationStore) GetListByName(ctx context.Context, namespaceId int, name string) ([]*entry.Application, error) {
	return a.ListQuery(ctx, "`namespace` = ? and like ?", []interface{}{namespaceId, "%" + name + "%"}, "create_time asc")
}

func newApplicationStore(db IDB) IApplicationStore {
	return &applicationStore{BaseStore: CreateStore[entry.Application](db)}
}
