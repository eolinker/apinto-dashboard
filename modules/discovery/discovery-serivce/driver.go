package discovery_serivce

import (
	driver_manager "github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
)

const (
	DriverStatic = "static"
	DriverConsul = "consul"
	DriverNacos  = "nacos"
	DriverEureka = "eureka"
)

type discoveryDriver struct {
	*driver_manager.DriverManager[discovery.IDiscoveryDriver]
}

func (d *discoveryDriver) List() []*driver_manager.DriverInfo {
	list := make([]*driver_manager.DriverInfo, 0)
	for name, value := range d.Drivers() {
		list = append(list, &driver_manager.DriverInfo{
			Name:   name,
			Render: value.Render(),
		})
	}
	return list
}

func newDiscoveryDriverManager() discovery.IDiscoveryDriverManager {
	return &discoveryDriver{DriverManager: driver_manager.CreateDriverManager[discovery.IDiscoveryDriver]()}
}
