package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
)

var (
	_ IClusterStore = (*clusterStore)(nil)
)

type IClusterStore interface {
	IBaseStore[cluster_entry.Cluster]
	GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_entry.Cluster, error)
	GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_entry.Cluster, error)
	GetAllByNamespaceId(ctx context.Context, namespace int) ([]*cluster_entry.Cluster, error)
	GetAll(ctx context.Context) ([]*cluster_entry.Cluster, error)
}

type clusterStore struct {
	*BaseStore[cluster_entry.Cluster]
}

func newClusterStore(db IDB) IClusterStore {
	return &clusterStore{BaseStore: CreateStore[cluster_entry.Cluster](db)}
}

func (c *clusterStore) GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_entry.Cluster, error) {
	return c.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (c *clusterStore) GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `name` in (?)", []interface{}{namespaceId, names}, "")
}

func (c *clusterStore) GetAllByNamespaceId(ctx context.Context, namespace int) ([]*cluster_entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ?", []interface{}{namespace}, "id desc")
}

func (c *clusterStore) GetAll(ctx context.Context) ([]*cluster_entry.Cluster, error) {
	list := make([]*cluster_entry.Cluster, 0)

	err := c.DB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
