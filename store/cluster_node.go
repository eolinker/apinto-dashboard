package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var (
	_ IClusterNodeStore = (*clusterNodeStore)(nil)
)

type IClusterNodeStore interface {
	IBaseStore[entry.ClusterNode]
	GetAllByClusterIds(ctx context.Context, clusterIds ...int) ([]*entry.ClusterNode, error)
	UpdateNodes(ctx context.Context, clusterId int, nodes []*entry.ClusterNode) error
}

type clusterNodeStore struct {
	*baseStore[entry.ClusterNode]
}

func newClusterNodeStore(db IDB) IClusterNodeStore {
	return &clusterNodeStore{baseStore: createStore[entry.ClusterNode](db)}
}

func (c *clusterNodeStore) GetAllByClusterIds(ctx context.Context, clusterIds ...int) ([]*entry.ClusterNode, error) {
	return c.ListQuery(ctx, "cluster in (?)", []interface{}{clusterIds}, "")
}

// UpdateNodes 重置节点信息 先把原来的删除
func (c *clusterNodeStore) UpdateNodes(ctx context.Context, clusterId int, nodes []*entry.ClusterNode) error {
	//先删除
	if err := c.DB(ctx).Exec("delete from `cluster_node` where cluster = ?", clusterId).Error; err != nil {
		return err
	}
	//再添加
	return c.DB(ctx).Create(nodes).Error
}
