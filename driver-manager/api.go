package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

type IAPIDriverManager interface {
	IDriverManager[driver.IAPIDriver]
	List() []*APIDriverInfo
}

type APIDriverInfo struct {
	Name string
}

type apiDriver struct {
	*driverManager[driver.IAPIDriver]
}

func (d *apiDriver) List() []*APIDriverInfo {
	list := make([]*APIDriverInfo, 0)
	for name, _ := range d.drivers {
		list = append(list, &APIDriverInfo{
			Name: name,
		})
	}
	return list
}

func newAPIDriverManager() IAPIDriverManager {
	return &apiDriver{driverManager: createDriverManager[driver.IAPIDriver]()}
}
