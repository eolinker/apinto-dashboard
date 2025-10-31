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
	oauth2 := driver.CreateOauth2()
	paraHmac := driver.CreateParaHmac()
	openidConnectJWT := driver.CreateOpenidConnectJWT()
	authDriverManager := newAuthDriverManager()
	authDriverManager.RegisterDriver(DriverApikey, apikey)
	authDriverManager.RegisterDriver(DriverAKsK, aksk)
	authDriverManager.RegisterDriver(DriverJwt, jwt)
	authDriverManager.RegisterDriver(DriverBasic, basic)
	authDriverManager.RegisterDriver(DriverOauth2, oauth2)
	authDriverManager.RegisterDriver(DriverOpneidConnectJWT, openidConnectJWT)
	authDriverManager.RegisterDriver(DriverParaHmac, paraHmac)

	bean.Injection(&authDriverManager)

	application := newApplicationService()

	bean.Injection(&application)
}
