package driver

import (
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/monitor"
)

const (
	MonitorInflux2 = "influxV2"
)

type monitorSourceManager struct {
	*driver.DriverManager[monitor.IMonitorSourceDriver]
}

func (a *monitorSourceManager) List() []string {
	return []string{MonitorInflux2}
}

func NewMonitorManager() monitor.IMonitorSourceManager {
	return &monitorSourceManager{DriverManager: driver.CreateDriverManager[monitor.IMonitorSourceDriver]()}
}
