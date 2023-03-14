package cluster_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IClusterConfigStore = (*clusterConfigStore)(nil)
)

type IClusterConfigStore interface {
	store.IBaseStore[cluster_entry.ClusterConfig]
	GetConfigByTypeByCluster(ctx context.Context, clusterID int, configType string) (*cluster_entry.ClusterConfig, error)
	GetConfigsByClusterID(ctx context.Context, clusterID int) ([]*cluster_entry.ClusterConfig, error)
}

type clusterConfigStore struct {
	*store.BaseStore[cluster_entry.ClusterConfig]
}

func newClusterConfigStore(db store.IDB) IClusterConfigStore {
	return &clusterConfigStore{BaseStore: store.CreateStore[cluster_entry.ClusterConfig](db)}
}

func (c *clusterConfigStore) GetConfigByTypeByCluster(ctx context.Context, clusterID int, configType string) (*cluster_entry.ClusterConfig, error) {
	config := new(cluster_entry.ClusterConfig)
	err := c.DB(ctx).Where("`cluster` = ? and `type` = ?", clusterID, configType).Take(config).Error
	return config, err
}

func (c *clusterConfigStore) GetConfigsByClusterID(ctx context.Context, clusterID int) ([]*cluster_entry.ClusterConfig, error) {
	return c.List(ctx, map[string]interface{}{
		"cluster": clusterID,
	})
}
