package service

import (
	"github.com/eolinker/apinto-dashboard/modules/monitor/driver"
	"github.com/eolinker/eosc/common/bean"
)

func init() {

	monitor := newMonitorService()
	monitorStatisticsService := newMonitorStatistics()
	monitorStatisticsCacheService := newMonitorStatisticsCache()
	monitorManager := driver.NewMonitorManager()
	apiHttp := driver.CreateMonitorInfluxV2(driver.MonitorInflux2)
	monitorManager.RegisterDriver(driver.MonitorInflux2, apiHttp)
	bean.Injection(&monitorManager)

	bean.Injection(&monitor)
	bean.Injection(&monitorStatisticsService)
	bean.Injection(&monitorStatisticsCacheService)
}
