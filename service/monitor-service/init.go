package monitor_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	iWarnStrategyService := newWarnStrategyService()
	iWarnHistoryService := newWarnHistoryService()
	iMonStatistics := newMonitorStatistics()
	iMonStatisticsCache := newMonitorStatisticsCache()
	monitor := newMonitorService()

	bean.Injection(&monitor)
	bean.Injection(&iWarnStrategyService)
	bean.Injection(&iWarnHistoryService)

	bean.Injection(&iMonStatistics)
	bean.Injection(&iMonStatisticsCache)
}
