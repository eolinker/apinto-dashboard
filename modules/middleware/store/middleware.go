package store

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/middleware/entry"

	"github.com/eolinker/apinto-dashboard/store"
)

type IMiddlewareStore interface {
	store.IBaseStore[entry.Middleware]
	GetByUUID(ctx context.Context, uuid string) (*entry.Middleware, error)
}

type middlewareStore struct {
	*store.BaseStore[entry.Middleware]
}

func (m *middlewareStore) GetByUUID(ctx context.Context, uuid string) (*entry.Middleware, error) {
	return m.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func newStore(db store.IDB) IMiddlewareStore {
	return &middlewareStore{BaseStore: store.CreateStore[entry.Middleware](db)}
}
