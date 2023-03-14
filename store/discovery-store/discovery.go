package discovery_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/discovery-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDiscoveryStore interface {
	store.IBaseStore[discovery_entry.Discovery]
	GetList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_entry.Discovery, error)
	GetByName(ctx context.Context, namespaceID int, discoveryName string) (*discovery_entry.Discovery, error)
}

type discoveryStore struct {
	*store.BaseStore[discovery_entry.Discovery]
}

func newDiscoveryStore(db store.IDB) IDiscoveryStore {
	return &discoveryStore{BaseStore: store.CreateStore[discovery_entry.Discovery](db)}
}

func (d *discoveryStore) GetList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_entry.Discovery, error) {
	if searchName != "" {
		return d.ListQuery(ctx, `namespace = ? and name like ?`, []interface{}{namespaceID, "%" + searchName + "%"}, "update_time DESC")
	} else {
		return d.ListQuery(ctx, `namespace = ?`, []interface{}{namespaceID}, "update_time DESC")
	}
}

func (d *discoveryStore) GetByName(ctx context.Context, namespaceID int, discoveryName string) (*discovery_entry.Discovery, error) {
	return d.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceID, discoveryName}, "")
}
