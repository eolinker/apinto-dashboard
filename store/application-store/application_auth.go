package application_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationAuthStore interface {
	store.IBaseStore[application_entry.ApplicationAuth]
	GetByUUID(ctx context.Context, uuid string) (*application_entry.ApplicationAuth, error)
	GetListByApplication(ctx context.Context, application int) ([]*application_entry.ApplicationAuth, error)
}

type applicationAuthStore struct {
	*store.BaseStore[application_entry.ApplicationAuth]
}

func newApplicationAuthStore(db store.IDB) IApplicationAuthStore {
	return &applicationAuthStore{BaseStore: store.CreateStore[application_entry.ApplicationAuth](db)}
}

func (a *applicationAuthStore) GetByUUID(ctx context.Context, uuid string) (*application_entry.ApplicationAuth, error) {
	return a.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (a *applicationAuthStore) GetListByApplication(ctx context.Context, application int) ([]*application_entry.ApplicationAuth, error) {
	return a.ListQuery(ctx, "`application`", []interface{}{application}, "`update_time` asc")
}
