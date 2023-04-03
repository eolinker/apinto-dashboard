package cluster_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		cluster := newClusterStore(db)
		clusterHistory := newClusterHistoryStore(db)
		clusterCertificate := newClusterCertificateStore(db)
		clusterNode := newClusterNodeStore(db)
		clusterRuntime := newClusterRuntimeStore(db)
		clusterConfig := newClusterConfigStore(db)
		clusterConfigRuntime := newClusterConfigRuntimeStore(db)

		bean.Injection(&cluster)
		bean.Injection(&clusterHistory)
		bean.Injection(&clusterCertificate)

		bean.Injection(&clusterNode)
		bean.Injection(&clusterRuntime)
		bean.Injection(&clusterConfig)
		bean.Injection(&clusterConfigRuntime)
	})
}
