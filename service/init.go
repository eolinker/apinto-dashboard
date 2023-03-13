package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	apintoClient := newApintoClientService()
	cluster := newClusterService()
	clusterCertificate := newClusterCertificateService()
	clusterVariable := newClusterVariableService()
	namespace := newNamespaceService()
	clusterNode := newClusterNodeService()
	clusterConfig := newClusterConfigService()
	globalVariable := newGlobalVariableService()
	discovery := newDiscoveryService()

	group := newCommonGroupService()
	application := newApplicationService()
	applicationAuth := newApplicationAuth()
	random := newRandomService()

	commonStrategy := newStrategyCommonService()
	lockAsynService := NewAsynLockService()
	lockSyncService := NewSyncLockService()

	userInfo := newUserInfoService()
	auditLog := newAuditLogService()
	bussinessAuth := newBussinessAuthService()
	externalAPP := newExternalApplicationService()
	monitor := newMonitorService()

	//openAPI
	apiOpenAPI := newAPIOpenAPIService()
	iMonStatistics := newMonitorStatistics()
	iMonStatisticsCache := newMonitorStatisticsCache()

	iNoticeChannelService := newNoticeChannelService()

	iWarnStrategyService := newWarnStrategyService()
	iWarnHistoryService := newWarnHistoryService()

	bean.Injection(&cluster)
	bean.Injection(&clusterCertificate)
	bean.Injection(&clusterVariable)
	bean.Injection(&clusterNode)
	bean.Injection(&clusterConfig)
	bean.Injection(&globalVariable)
	bean.Injection(&discovery)

	bean.Injection(&group)
	bean.Injection(&namespace)
	bean.Injection(&apintoClient)
	bean.Injection(&application)
	bean.Injection(&applicationAuth)
	bean.Injection(&random)
	bean.Injection(&lockAsynService)
	bean.Injection(&lockSyncService)

	bean.Injection(&commonStrategy)

	bean.Injection(&userInfo)
	bean.Injection(&auditLog)
	bean.Injection(&bussinessAuth)
	bean.Injection(&externalAPP)
	bean.Injection(&monitor)

	//openAPI
	bean.Injection(&apiOpenAPI)

	resetOnline := newResetOnline()
	bean.Injection(&resetOnline)
	bean.Injection(&iMonStatistics)
	bean.Injection(&iMonStatisticsCache)

	bean.Injection(&iNoticeChannelService)
	bean.Injection(&iWarnStrategyService)
	bean.Injection(&iWarnHistoryService)
}
