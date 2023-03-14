package application_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	application := newApplicationService()
	applicationAuth := newApplicationAuth()
	bean.Injection(&application)
	bean.Injection(&applicationAuth)
}
