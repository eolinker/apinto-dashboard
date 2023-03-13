package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var (
	_ IClusterRuntimeStore = (*clusterRuntimeStore)(nil)
)

type IClusterRuntimeStore interface {
	IBaseStore[entry.Runtime]
	DeleteByClusterID(ctx context.Context, clusterId int) error
}

type clusterRuntimeStore struct {
	*baseStore[entry.Runtime]
}

func newClusterRuntimeStore(db IDB) IClusterRuntimeStore {
	return &clusterRuntimeStore{baseStore: createStore[entry.Runtime](db)}
}

func (c *clusterRuntimeStore) DeleteByClusterID(ctx context.Context, clusterId int) error {
	delMap := map[string]interface{}{`cluster`: clusterId}
	_, err := c.DeleteWhere(ctx, delMap)
	return err
}
