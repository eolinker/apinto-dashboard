package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IDiscoveryStore interface {
	IBaseStore[entry.Discovery]
	GetList(ctx context.Context, namespaceID int, searchName string) ([]*entry.Discovery, error)
	GetByName(ctx context.Context, namespaceID int, discoveryName string) (*entry.Discovery, error)
}

type discoveryStore struct {
	*baseStore[entry.Discovery]
}

func newDiscoveryStore(db IDB) IDiscoveryStore {
	return &discoveryStore{baseStore: createStore[entry.Discovery](db)}
}

func (d *discoveryStore) GetList(ctx context.Context, namespaceID int, searchName string) ([]*entry.Discovery, error) {
	if searchName != "" {
		return d.ListQuery(ctx, `namespace = ? and name like ?`, []interface{}{namespaceID, "%" + searchName + "%"}, "update_time DESC")
	} else {
		return d.ListQuery(ctx, `namespace = ?`, []interface{}{namespaceID}, "update_time DESC")
	}
}

func (d *discoveryStore) GetByName(ctx context.Context, namespaceID int, discoveryName string) (*entry.Discovery, error) {
	return d.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceID, discoveryName}, "")
}
