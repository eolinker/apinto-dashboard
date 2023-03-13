package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

type IDiscoveryDriverManager interface {
	IDriverManager[driver.IDiscoveryDriver]
	List() []*DriverInfo
}

type discoveryDriver struct {
	*driverManager[driver.IDiscoveryDriver]
}

func (d *discoveryDriver) List() []*DriverInfo {
	list := make([]*DriverInfo, 0)
	for name, value := range d.drivers {
		list = append(list, &DriverInfo{
			Name:   name,
			Render: value.Render(),
		})
	}
	return list
}

func newDiscoveryDriverManager() IDiscoveryDriverManager {
	return &discoveryDriver{driverManager: createDriverManager[driver.IDiscoveryDriver]()}
}
