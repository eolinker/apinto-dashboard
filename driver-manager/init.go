package driver_manager

import (
	"github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	discoveryDriverManager := newDiscoveryDriverManager()

	consul := driver.CreateConsul("consul")
	nacos := driver.CreateNacos("nacos")
	eureka := driver.CreateEureka("eureka")
	staticService := driver.CreateStaticEnum("http")

	discoveryDriverManager.RegisterDriver(driver.DriverConsul, consul)
	discoveryDriverManager.RegisterDriver(driver.DriverNacos, nacos)
	discoveryDriverManager.RegisterDriver(driver.DriverEureka, eureka)

	apikey := driver.CreateApikey()
	aksk := driver.CreateAkSk()
	jwt := driver.CreateJwt()
	basic := driver.CreateBasic()

	authDriverManager := newAuthDriverManager()
	authDriverManager.RegisterDriver(driver.DriverApikey, apikey)
	authDriverManager.RegisterDriver(driver.DriverAKsK, aksk)
	authDriverManager.RegisterDriver(driver.DriverJwt, jwt)
	authDriverManager.RegisterDriver(driver.DriverBasic, basic)

	apiDriverManager := newAPIDriverManager()
	apiHttp := driver.CreateAPIHttp("http")
	apiDriverManager.RegisterDriver(driver.DriverApiHTTP, apiHttp)

	//strategyDriverManager := newStrategyDriverManager()
	//traffic := driver.CreateTraffic("limiting")
	//strategyDriverManager.registerDriver(enum.StrategyTraffic, traffic)

	clConfigDriverManager := newCLConfigDriverManager()
	redisDriver := driver.CreateRedis("redis")
	influxv2Driver := driver.CreateInfluxV2("influxdbv2")
	clConfigDriverManager.RegisterDriver(enum.CLConfigRedis, redisDriver)
	clConfigDriverManager.RegisterDriver(enum.CLConfigInfluxV2, influxv2Driver)

	//同步api文档格式管理器
	apiSyncFormatManager := newAPISyncFormatManager()
	openAPI3Driver := driver.CreateOpenAPI3(openAPI3)
	openAPI2Driver := driver.CreateOpenAPI2(swagger2)
	apiSyncFormatManager.RegisterDriver(openAPI3, openAPI3Driver)
	apiSyncFormatManager.RegisterDriver(swagger2, openAPI2Driver)

	//监控数据源管理器
	monitorSources := newMonitorManager()
	monInflux2 := driver.CreateMonitorInfluxV2(monitorInflux2)
	monitorSources.RegisterDriver(monitorInflux2, monInflux2)

	//渠道通知管理器
	channelDriverManager := newNoticeChannelDriverManager()

	bean.Injection(&discoveryDriverManager)
	bean.Injection(&authDriverManager)
	bean.Injection(&staticService)
	bean.Injection(&apiDriverManager)
	bean.Injection(&clConfigDriverManager)
	bean.Injection(&apiSyncFormatManager)
	bean.Injection(&monitorSources)
	bean.Injection(&channelDriverManager)
}
