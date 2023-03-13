package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

const (
	monitorInflux2 = "influxV2"
)

type IMonitorSourceManager interface {
	IDriverManager[driver.IMonitorSourceDriver]
	List() []string
}

type monitorSourceManager struct {
	*driverManager[driver.IMonitorSourceDriver]
}

func (a *monitorSourceManager) List() []string {
	return []string{monitorInflux2}
}

func newMonitorManager() IMonitorSourceManager {
	return &monitorSourceManager{driverManager: createDriverManager[driver.IMonitorSourceDriver]()}
}
