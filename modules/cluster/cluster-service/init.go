package cluster_service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {

	//clConfigDriverManager := newCLConfigDriverManager()
	//redisDriver := driver2.CreateRedis("redis")
	//clConfigDriverManager.RegisterDriver(cluster.CLConfigRedis, redisDriver)
	//
	//bean.Injection(&clConfigDriverManager)

	iClusterService := newClusterService()
	clusterCertificate := newClusterCertificateService()
	clusterNode := newClusterNodeService()
	//clusterConfig := newClusterConfigService()

	apintoClient := newApintoClientService()

	bean.Injection(&apintoClient)

	bean.Injection(&iClusterService)
	bean.Injection(&clusterCertificate)
	bean.Injection(&clusterNode)
	//bean.Injection(&clusterConfig)
	nodeCache := newINodeCache()
	bean.Injection(&nodeCache)

}
