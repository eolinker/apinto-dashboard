package api_service

import (
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/api"
)

const (
	DriverApiHTTP = "http"
)

type apiDriver struct {
	*driver.DriverManager[api.IAPIDriver]
}

func (d *apiDriver) List() []*api.APIDriverInfo {
	list := make([]*api.APIDriverInfo, 0)
	for name, _ := range d.Drivers() {
		list = append(list, &api.APIDriverInfo{
			Name: name,
		})
	}
	return list
}

func newAPIDriverManager() api.IAPIDriverManager {
	return &apiDriver{DriverManager: driver.CreateDriverManager[api.IAPIDriver]()}
}
