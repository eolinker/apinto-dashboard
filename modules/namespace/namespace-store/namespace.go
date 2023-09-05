package namespace_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"time"
)

type INamespaceStore interface {
	store.IBaseStore[namespace_entry.Namespace]
	GetByName(ctx context.Context, name string) (*namespace_entry.Namespace, error)
	GetAll(ctx context.Context) ([]*namespace_entry.Namespace, error)
}
type INamespaceCache cache.IRedisCache[namespace_entry.Namespace, string]

type namespaceStore struct {
	*store.BaseStore[namespace_entry.Namespace]
}

func newNamespaceStore(db store.IDB) INamespaceStore {

	s := &namespaceStore{BaseStore: store.CreateStore[namespace_entry.Namespace](db)}
	s.Save(context.Background(), &namespace_entry.Namespace{
		Id:         1,
		Name:       "default",
		CreateTime: time.Now(),
	})
	return s
}

func (n *namespaceStore) GetByName(ctx context.Context, name string) (*namespace_entry.Namespace, error) {
	return n.FirstQuery(ctx, "`name` = ?", []interface{}{name}, "")
}

func (n *namespaceStore) GetAll(ctx context.Context) ([]*namespace_entry.Namespace, error) {
	return n.List(ctx, map[string]interface{}{})
}
