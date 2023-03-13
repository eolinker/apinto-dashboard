package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var (
	_ IClusterConfigStore = (*clusterConfigStore)(nil)
)

type IClusterConfigStore interface {
	IBaseStore[entry.ClusterConfig]
	GetConfigByTypeByCluster(ctx context.Context, clusterID int, configType string) (*entry.ClusterConfig, error)
	GetConfigsByClusterID(ctx context.Context, clusterID int) ([]*entry.ClusterConfig, error)
}

type clusterConfigStore struct {
	*BaseStore[entry.ClusterConfig]
}

func newClusterConfigStore(db IDB) IClusterConfigStore {
	return &clusterConfigStore{BaseStore: CreateStore[entry.ClusterConfig](db)}
}

func (c *clusterConfigStore) GetConfigByTypeByCluster(ctx context.Context, clusterID int, configType string) (*entry.ClusterConfig, error) {
	config := new(entry.ClusterConfig)
	err := c.DB(ctx).Where("`cluster` = ? and `type` = ?", clusterID, configType).Take(config).Error
	return config, err
}

func (c *clusterConfigStore) GetConfigsByClusterID(ctx context.Context, clusterID int) ([]*entry.ClusterConfig, error) {
	return c.List(ctx, map[string]interface{}{
		"cluster": clusterID,
	})
}
