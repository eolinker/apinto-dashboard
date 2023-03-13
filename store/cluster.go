package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var (
	_ IClusterStore = (*clusterStore)(nil)
)

type IClusterStore interface {
	IBaseStore[entry.Cluster]
	GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*entry.Cluster, error)
	GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*entry.Cluster, error)
	GetAllByNamespaceId(ctx context.Context, namespace int) ([]*entry.Cluster, error)
	GetAll(ctx context.Context) ([]*entry.Cluster, error)
}

type clusterStore struct {
	*baseStore[entry.Cluster]
}

func newClusterStore(db IDB) IClusterStore {
	return &clusterStore{baseStore: createStore[entry.Cluster](db)}
}

func (c *clusterStore) GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*entry.Cluster, error) {
	return c.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (c *clusterStore) GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `name` in (?)", []interface{}{namespaceId, names}, "")
}

func (c *clusterStore) GetAllByNamespaceId(ctx context.Context, namespace int) ([]*entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ?", []interface{}{namespace}, "id desc")
}

func (c *clusterStore) GetAll(ctx context.Context) ([]*entry.Cluster, error) {
	list := make([]*entry.Cluster, 0)

	err := c.DB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
