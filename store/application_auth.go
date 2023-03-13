package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthStore interface {
	IBaseStore[entry.ApplicationAuth]
	GetByUUID(ctx context.Context, uuid string) (*entry.ApplicationAuth, error)
	GetListByApplication(ctx context.Context, application int) ([]*entry.ApplicationAuth, error)
}

type applicationAuthStore struct {
	*baseStore[entry.ApplicationAuth]
}

func newApplicationAuthStore(db IDB) IApplicationAuthStore {
	return &applicationAuthStore{baseStore: createStore[entry.ApplicationAuth](db)}
}

func (a *applicationAuthStore) GetByUUID(ctx context.Context, uuid string) (*entry.ApplicationAuth, error) {
	return a.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (a *applicationAuthStore) GetListByApplication(ctx context.Context, application int) ([]*entry.ApplicationAuth, error) {
	return a.ListQuery(ctx, "`application`", []interface{}{application}, "`update_time` asc")
}
