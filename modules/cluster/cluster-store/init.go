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
		gmCertificate := newGMClusterCertificateStore(db)
		clusterNode := newClusterNodeStore(db)
		//clusterConfig := newClusterConfigStore(db)
		//bean.Injection(&clusterConfig)

		bean.Injection(&cluster)
		bean.Injection(&clusterHistory)
		bean.Injection(&clusterCertificate)
		bean.Injection(&gmCertificate)

		bean.Injection(&clusterNode)
		clusterRuntime := newClusterRuntimeStore(db)
		bean.Injection(&clusterRuntime)
		//clusterConfigRuntime := newClusterConfigRuntimeStore(db)
		//bean.Injection(&clusterConfigRuntime)
	})
}
