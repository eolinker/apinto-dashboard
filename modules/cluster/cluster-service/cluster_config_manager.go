package cluster_service

import (
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
)

type clConfigManager struct {
	*driver.DriverManager[cluster.ICLConfigDriver]
}

func (d *clConfigManager) List() []*driver.DriverInfo {
	list := make([]*driver.DriverInfo, 0)
	for name := range d.Drivers() {
		list = append(list, &driver.DriverInfo{
			Name: name,
		})
	}
	return list
}

func newCLConfigDriverManager() cluster.ICLConfigDriverManager {
	return &clConfigManager{DriverManager: driver.CreateDriverManager[cluster.ICLConfigDriver]()}
}
