package application_service

import (
	"github.com/eolinker/apinto-dashboard/modules/application/driver"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	apikey := driver.CreateApikey()
	aksk := driver.CreateAkSk()
	jwt := driver.CreateJwt()
	basic := driver.CreateBasic()

	authDriverManager := newAuthDriverManager()
	authDriverManager.RegisterDriver(DriverApikey, apikey)
	authDriverManager.RegisterDriver(DriverAKsK, aksk)
	authDriverManager.RegisterDriver(DriverJwt, jwt)
	authDriverManager.RegisterDriver(DriverBasic, basic)

	bean.Injection(&authDriverManager)

	application := newApplicationService()
	applicationAuth := newApplicationAuth()
	bean.Injection(&application)
	bean.Injection(&applicationAuth)
}
