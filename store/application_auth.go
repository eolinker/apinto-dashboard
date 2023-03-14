package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
)

type IApplicationAuthStore interface {
	IBaseStore[application_entry.ApplicationAuth]
	GetByUUID(ctx context.Context, uuid string) (*application_entry.ApplicationAuth, error)
	GetListByApplication(ctx context.Context, application int) ([]*application_entry.ApplicationAuth, error)
}

type applicationAuthStore struct {
	*BaseStore[application_entry.ApplicationAuth]
}

func newApplicationAuthStore(db IDB) IApplicationAuthStore {
	return &applicationAuthStore{BaseStore: CreateStore[application_entry.ApplicationAuth](db)}
}

func (a *applicationAuthStore) GetByUUID(ctx context.Context, uuid string) (*application_entry.ApplicationAuth, error) {
	return a.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (a *applicationAuthStore) GetListByApplication(ctx context.Context, application int) ([]*application_entry.ApplicationAuth, error) {
	return a.ListQuery(ctx, "`application`", []interface{}{application}, "`update_time` asc")
}
