package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
)

var (
	_ IClusterCertificateStore = (*clusterCertificateStore)(nil)
)

type IClusterCertificateStore interface {
	IBaseStore[cluster_entry.ClusterCertificate]
	QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterCertificate, error)
}

type clusterCertificateStore struct {
	*BaseStore[cluster_entry.ClusterCertificate]
}

func newClusterCertificateStore(db IDB) IClusterCertificateStore {
	return &clusterCertificateStore{BaseStore: CreateStore[cluster_entry.ClusterCertificate](db)}
}

func (c *clusterCertificateStore) QueryListByClusterId(ctx context.Context, clusterId int) ([]*cluster_entry.ClusterCertificate, error) {
	return c.ListQuery(ctx, "`cluster` = ?", []interface{}{clusterId}, "")
}
