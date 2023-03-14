package bussiness_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	bussinessAuth := newBussinessAuthService()
	bean.Injection(&bussinessAuth)

}
