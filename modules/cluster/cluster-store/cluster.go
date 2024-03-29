package cluster_store

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IClusterStore = (*clusterStore)(nil)
)

type IClusterStore interface {
	store.IBaseStore[cluster_entry.Cluster]
	GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_entry.Cluster, error)
	GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_entry.Cluster, error)
	GetByNamespaceByUUIDs(ctx context.Context, namespaceId int, UUIDs []string) ([]*cluster_entry.Cluster, error)
	GetAllByNamespaceId(ctx context.Context, namespace int) ([]*cluster_entry.Cluster, error)
	GetAll(ctx context.Context) ([]*cluster_entry.Cluster, error)
	ClusterCount(ctx context.Context, params map[string]interface{}) (int64, error)
}

type clusterStore struct {
	*store.BaseStore[cluster_entry.Cluster]
}

func newClusterStore(db store.IDB) IClusterStore {
	return &clusterStore{BaseStore: store.CreateStore[cluster_entry.Cluster](db)}
}

func (c *clusterStore) ClusterCount(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	err := c.DB(ctx).Where(params).Model(cluster_entry.Cluster{}).Count(&count).Error
	return count, err
}

func (c *clusterStore) GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_entry.Cluster, error) {
	return c.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (c *clusterStore) GetByNamespaceByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `name` in (?)", []interface{}{namespaceId, names}, "")
}

func (c *clusterStore) GetByNamespaceByUUIDs(ctx context.Context, namespaceId int, UUIDs []string) ([]*cluster_entry.Cluster, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `uuid` in (?)", []interface{}{namespaceId, UUIDs}, "")
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
