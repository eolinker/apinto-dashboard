package cluster_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	cluster := newClusterService()
	clusterCertificate := newClusterCertificateService()
	clusterNode := newClusterNodeService()
	clusterConfig := newClusterConfigService()

	apintoClient := newApintoClientService()

	bean.Injection(&apintoClient)

	bean.Injection(&cluster)
	bean.Injection(&clusterCertificate)
	bean.Injection(&clusterNode)
	bean.Injection(&clusterConfig)
}
