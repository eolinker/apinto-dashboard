package openapp_service

import (
	driver2 "github.com/eolinker/apinto-dashboard/modules/openapi/driver"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	//同步api文档格式管理器
	apiSyncFormatManager := newAPISyncFormatManager()
	openAPI3Driver := driver2.CreateOpenAPI3(openAPI3)
	openAPI2Driver := driver2.CreateOpenAPI2(swagger2)
	apiSyncFormatManager.RegisterDriver(openAPI3, openAPI3Driver)
	apiSyncFormatManager.RegisterDriver(swagger2, openAPI2Driver)
	bean.Injection(&apiSyncFormatManager)

	externalAPP := newExternalApplicationService()

	bean.Injection(&externalAPP)
}
