package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

var (
	_ IClusterCertificateStore = (*clusterCertificateStore)(nil)
)

type IClusterCertificateStore interface {
	IBaseStore[entry.ClusterCertificate]
	QueryListByClusterId(ctx context.Context, clusterId int) ([]*entry.ClusterCertificate, error)
}

type clusterCertificateStore struct {
	*baseStore[entry.ClusterCertificate]
}

func newClusterCertificateStore(db IDB) IClusterCertificateStore {
	return &clusterCertificateStore{baseStore: createStore[entry.ClusterCertificate](db)}
}

func (c *clusterCertificateStore) QueryListByClusterId(ctx context.Context, clusterId int) ([]*entry.ClusterCertificate, error) {
	return c.ListQuery(ctx, "`cluster` = ?", []interface{}{clusterId}, "")
}
