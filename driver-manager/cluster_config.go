package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

// ICLConfigDriverManager 集群配置驱动管理器
type ICLConfigDriverManager interface {
	IDriverManager[driver.ICLConfigDriver]
	List() []*DriverInfo
}

type clConfigManager struct {
	*driverManager[driver.ICLConfigDriver]
}

func (d *clConfigManager) List() []*DriverInfo {
	list := make([]*DriverInfo, 0)
	for name := range d.drivers {
		list = append(list, &DriverInfo{
			Name: name,
		})
	}
	return list
}

func newCLConfigDriverManager() ICLConfigDriverManager {
	return &clConfigManager{driverManager: createDriverManager[driver.ICLConfigDriver]()}
}
