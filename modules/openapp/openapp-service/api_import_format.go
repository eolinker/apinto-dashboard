package openapp_service

import (
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/openapi"
)

const (
	openAPI3 = "OpenAPI3.0"
	swagger2 = "Swagger2.0"
)

type apiSyncFormatDriver struct {
	*driver.DriverManager[openapi.IAPISyncFormatDriver]
}

func (a *apiSyncFormatDriver) List() []string {
	return []string{swagger2, openAPI3}
}

func newAPISyncFormatManager() openapi.IAPISyncFormatManager {
	return &apiSyncFormatDriver{DriverManager: driver.CreateDriverManager[openapi.IAPISyncFormatDriver]()}
}
