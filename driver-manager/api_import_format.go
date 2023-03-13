package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

const (
	openAPI3 = "OpenAPI3.0"
	swagger2 = "Swagger2.0"
)

type IAPISyncFormatManager interface {
	IDriverManager[driver.IAPISyncFormatDriver]
	List() []string
}

type apiSyncFormatDriver struct {
	*driverManager[driver.IAPISyncFormatDriver]
}

func (a *apiSyncFormatDriver) List() []string {
	return []string{swagger2, openAPI3}
}

func newAPISyncFormatManager() IAPISyncFormatManager {
	return &apiSyncFormatDriver{driverManager: createDriverManager[driver.IAPISyncFormatDriver]()}
}
