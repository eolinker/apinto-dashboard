package cluster_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/runtime-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IClusterRuntimeStore = (*clusterRuntimeStore)(nil)
)

type IClusterRuntimeStore interface {
	store.IBaseStore[runtime_entry.Runtime]
	DeleteByClusterID(ctx context.Context, clusterId int) error
}

type clusterRuntimeStore struct {
	*store.BaseStore[runtime_entry.Runtime]
}

func newClusterRuntimeStore(db store.IDB) IClusterRuntimeStore {
	return &clusterRuntimeStore{BaseStore: store.CreateStore[runtime_entry.Runtime](db)}
}

func (c *clusterRuntimeStore) DeleteByClusterID(ctx context.Context, clusterId int) error {
	delMap := map[string]interface{}{`cluster`: clusterId}
	_, err := c.DeleteWhere(ctx, delMap)
	return err
}
