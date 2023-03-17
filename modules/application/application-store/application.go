package application_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var _ IApplicationStore = (*applicationStore)(nil)

type IApplicationStore interface {
	store.IBaseStore[application_entry.Application]
	GetByIdStr(ctx context.Context, namespaceId int, idStr string) (*application_entry.Application, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*application_entry.Application, error)
	GetListPage(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*application_entry.Application, int, error)
	GetList(ctx context.Context, namespaceId int, uuids ...string) ([]*application_entry.Application, error)
	GetListByName(ctx context.Context, namespaceId int, name string) ([]*application_entry.Application, error)
}

type applicationStore struct {
	*store.BaseStore[application_entry.Application]
}

func (a *applicationStore) GetByIdStr(ctx context.Context, namespaceId int, idStr string) (*application_entry.Application, error) {
	return a.FirstQuery(ctx, "`namespace` = ? and `id_str` = ?", []interface{}{namespaceId, idStr}, "")
}

func (a *applicationStore) GetByName(ctx context.Context, namespaceId int, name string) (*application_entry.Application, error) {
	return a.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (a *applicationStore) GetListPage(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*application_entry.Application, int, error) {
	if queryName == "" {
		return a.ListPage(ctx, "`namespace` = ?", pageNum, pageSize, []interface{}{namespaceId}, "update_time desc")
	}
	return a.ListPage(ctx, "`namespace` = ? and `name` like ?", pageNum, pageSize, []interface{}{namespaceId, "%" + queryName + "%"}, "update_time desc")
}

func (a *applicationStore) GetList(ctx context.Context, namespaceId int, uuids ...string) ([]*application_entry.Application, error) {
	if len(uuids) > 0 {
		return a.ListQuery(ctx, "`namespace` = ? and `id_str` in (?)", []interface{}{namespaceId, uuids}, "create_time asc")
	}
	return a.ListQuery(ctx, "`namespace` = ?", []interface{}{namespaceId}, "create_time asc")
}

func (a *applicationStore) GetListByName(ctx context.Context, namespaceId int, name string) ([]*application_entry.Application, error) {
	return a.ListQuery(ctx, "`namespace` = ? and like ?", []interface{}{namespaceId, "%" + name + "%"}, "create_time asc")
}

func newApplicationStore(db store.IDB) IApplicationStore {
	return &applicationStore{BaseStore: store.CreateStore[application_entry.Application](db)}
}
