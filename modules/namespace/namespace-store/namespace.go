package namespace_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type INamespaceStore interface {
	store.IBaseStore[namespace_entry.Namespace]
	GetByName(ctx context.Context, name string) (*namespace_entry.Namespace, error)
	GetAll(ctx context.Context) ([]*namespace_entry.Namespace, error)
}

type namespaceStore struct {
	*store.BaseStore[namespace_entry.Namespace]
}

func newNamespaceStore(db store.IDB) INamespaceStore {
	return &namespaceStore{BaseStore: store.CreateStore[namespace_entry.Namespace](db)}
}

func (n *namespaceStore) GetByName(ctx context.Context, name string) (*namespace_entry.Namespace, error) {
	return n.FirstQuery(ctx, "`name` = ?", []interface{}{name}, "")
}

func (n *namespaceStore) GetAll(ctx context.Context) ([]*namespace_entry.Namespace, error) {
	return n.List(ctx, map[string]interface{}{})
}
