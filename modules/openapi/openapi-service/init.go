package openapi_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	apiOpenAPI := newAPIOpenAPIService()
	//openAPI
	bean.Injection(&apiOpenAPI)
}
