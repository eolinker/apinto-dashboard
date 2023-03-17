package discovery_serivce

import (
	driver2 "github.com/eolinker/apinto-dashboard/modules/discovery/driver"
	"github.com/eolinker/eosc/common/bean"
)

func init() {

	discoveryDriverManager := newDiscoveryDriverManager()

	consul := driver2.CreateConsul("consul")
	nacos := driver2.CreateNacos("nacos")
	eureka := driver2.CreateEureka("eureka")
	staticService := driver2.CreateStaticEnum("http")

	discoveryDriverManager.RegisterDriver(DriverConsul, consul)
	discoveryDriverManager.RegisterDriver(DriverNacos, nacos)
	discoveryDriverManager.RegisterDriver(DriverEureka, eureka)

	bean.Injection(&staticService)
	bean.Injection(&discoveryDriverManager)
	discovery := newDiscoveryService()

	bean.Injection(&discovery)
}
