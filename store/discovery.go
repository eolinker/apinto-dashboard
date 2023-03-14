package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/discovery-entry"
)

type IDiscoveryStore interface {
	IBaseStore[discovery_entry.Discovery]
	GetList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_entry.Discovery, error)
	GetByName(ctx context.Context, namespaceID int, discoveryName string) (*discovery_entry.Discovery, error)
}

type discoveryStore struct {
	*BaseStore[discovery_entry.Discovery]
}

func newDiscoveryStore(db IDB) IDiscoveryStore {
	return &discoveryStore{BaseStore: CreateStore[discovery_entry.Discovery](db)}
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
