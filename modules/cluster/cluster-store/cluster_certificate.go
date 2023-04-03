package cluster_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IClusterCertificateStore = (*clusterCertificateStore)(nil)
)

type IClusterCertificateStore interface {
	store.IBaseStore[cluster_entry.ClusterCertificate]
	QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterCertificate, error)
}

type clusterCertificateStore struct {
	*store.BaseStore[cluster_entry.ClusterCertificate]
}

func newClusterCertificateStore(db store.IDB) IClusterCertificateStore {
	return &clusterCertificateStore{BaseStore: store.CreateStore[cluster_entry.ClusterCertificate](db)}
}

func (c *clusterCertificateStore) QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterCertificate, error) {
	return c.ListQuery(ctx, "`cluster` = ?", []interface{}{clusterId}, "")
}
