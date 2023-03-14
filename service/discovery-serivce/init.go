package discovery_serivce

import "github.com/eolinker/eosc/common/bean"

func init() {
	discovery := newDiscoveryService()

	bean.Injection(&discovery)
}
