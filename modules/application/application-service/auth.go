package application_service

import (
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/application"
)

type authDriver struct {
	*driver.DriverManager[application.IAuthDriver]
}

func newAuthDriverManager() application.IAuthDriverManager {
	return &authDriver{DriverManager: driver.CreateDriverManager[application.IAuthDriver]()}
}

func (d *authDriver) List() []*driver.DriverInfo {
	list := make([]*driver.DriverInfo, 0)
	for name, value := range d.Drivers() {
		list = append(list, &driver.DriverInfo{
			Name:   name,
			Render: value.Render(),
		})
	}
	return list
}
