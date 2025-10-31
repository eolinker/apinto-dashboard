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

type IGMCertificateStore interface {
	store.IBaseStore[cluster_entry.ClusterGMCertificate]
	QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterGMCertificate, error)
}

type gmCertificateStore struct {
	*store.BaseStore[cluster_entry.ClusterGMCertificate]
}

func newGMClusterCertificateStore(db store.IDB) IGMCertificateStore {
	return &gmCertificateStore{BaseStore: store.CreateStore[cluster_entry.ClusterGMCertificate](db)}
}

func (c *gmCertificateStore) QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterGMCertificate, error) {
	return c.ListQuery(ctx, "`cluster` = ?", []interface{}{clusterId}, "")
}
