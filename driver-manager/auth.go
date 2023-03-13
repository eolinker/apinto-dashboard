package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

type IAuthDriverManager interface {
	IDriverManager[driver.IAuthDriver]
	List() []*DriverInfo
}

type authDriver struct {
	*driverManager[driver.IAuthDriver]
}

func newAuthDriverManager() IAuthDriverManager {
	return &authDriver{driverManager: createDriverManager[driver.IAuthDriver]()}
}

func (d *authDriver) List() []*DriverInfo {
	list := make([]*DriverInfo, 0)
	for name, value := range d.drivers {
		list = append(list, &DriverInfo{
			Name:   name,
			Render: value.Render(),
		})
	}
	return list
}
