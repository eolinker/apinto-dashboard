package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type INamespaceStore interface {
	IBaseStore[entry.Namespace]
	GetByName(ctx context.Context, name string) (*entry.Namespace, error)
	GetAll(ctx context.Context) ([]*entry.Namespace, error)
}

type namespaceStore struct {
	*baseStore[entry.Namespace]
}

func newNamespaceStore(db IDB) INamespaceStore {
	return &namespaceStore{baseStore: createStore[entry.Namespace](db)}
}

func (n *namespaceStore) GetByName(ctx context.Context, name string) (*entry.Namespace, error) {
	return n.FirstQuery(ctx, "`name` = ?", []interface{}{name}, "")
}

func (n *namespaceStore) GetAll(ctx context.Context) ([]*entry.Namespace, error) {
	return n.List(ctx, map[string]interface{}{})
}
