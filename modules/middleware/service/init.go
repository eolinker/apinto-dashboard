package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	middlewareService := newMiddlewareService()

	bean.Injection(&middlewareService)

}
