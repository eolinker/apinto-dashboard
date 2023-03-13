package store

import (
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

func InitStoreDB(idb *gorm.DB) {
	var db IDB = &myDB{db: idb}
	runHandler(db)
	innetHandler(db)

}
func innetHandler(db IDB) {
	cluster := newClusterStore(db)
	clusterHistory := newClusterHistoryStore(db)
	clusterCertificate := newClusterCertificateStore(db)
	clusterVariable := newClusterVariableStore(db)
	clusterNode := newClusterNodeStore(db)
	clusterRuntime := newClusterRuntimeStore(db)
	clusterConfig := newClusterConfigStore(db)
	clusterConfigRuntime := newClusterConfigRuntimeStore(db)

	globalVariable := newGlobalVariableStore(db)

	variableHistory := newVariableHistoryStore(db)
	variablePublishVersion := newVariablePublishVersionStore(db)
	variableRuntime := newVariableRuntimeStore(db)
	variablePublishHistory := newVariablePublishHistoryStore(db)

	discovery := newDiscoveryStore(db)

	namespace := newNamespaceStore(db)

	commonGroup := newCommonGroupStore(db)

	quote := newQuoteStore(db)

	discoveryVersionStore := newDiscoveryVersionStore(db)
	discoveryStatStore := newDiscoveryStatStore(db)
	discoveryRuntime := newDiscoveryRuntimeStore(db)
	discoveryHistory := newDiscoveryHistoryStore(db)

	application := newApplicationStore(db)
	applicationRuntime := newApplicationRuntimeStore(db)
	applicationStat := newApplicationStatStore(db)
	applicationVersion := newApplicationVersionStore(db)
	applicationAuth := newApplicationAuthStore(db)
	applicationAuthVersion := newApplicationAuthVersionStore(db)
	applicationAuthStat := newApplicationAuthStatStore(db)
	applicationAuthRuntimeStore := newApplicationAuthRuntimeStore(db)
	applicationAuthPublish := newApplicationAuthPublishStore(db)
	applicationHistory := newApplicationHistoryStore(db)
	applicationAuthHistory := newApplicationAuthHistoryStore(db)

	strategy := newStrategyStore(db)
	strategyStat := newStrategyStatStore(db)

	strategyVersion := newStrategyVersionStore(db)

	strategyHistory := newStrategyHistoryStore(db)

	userInfo := newUserInfoStore(db)
	role := newRoleStore(db)
	roleAccess := newRoleAccessStore(db)
	userRole := newUserRoleStore(db)
	roleAccessLog := newRoleAccessLogStore(db)

	auditLog := newAuditLogStore(db)
	systemInfo := newSystemInfoStore(db)
	externalAPP := newExternalApplicationStore(db)
	monitor := newMonitorStore(db)

	iNoticeChannelStore := newNoticeChannelStore(db)
	iNoticeChannelStatStore := newNoticeChannelStatStore(db)
	iNoticeChannelVersionStore := newNoticeChannelVersionStore(db)

	iWarnStrategyStore := newWarnStrategyIStore(db)
	iWarnHistoryStore := newWarnHistoryIStore(db)

	bean.Injection(&db)

	bean.Injection(&cluster)
	bean.Injection(&clusterHistory)
	bean.Injection(&clusterCertificate)
	bean.Injection(&clusterVariable)
	bean.Injection(&clusterNode)
	bean.Injection(&clusterRuntime)
	bean.Injection(&clusterConfig)
	bean.Injection(&clusterConfigRuntime)
	bean.Injection(&globalVariable)
	bean.Injection(&quote)
	bean.Injection(&variableHistory)
	bean.Injection(&variablePublishVersion)
	bean.Injection(&variableRuntime)
	bean.Injection(&variablePublishHistory)

	bean.Injection(&namespace)

	bean.Injection(&discovery)
	bean.Injection(&discoveryVersionStore)
	bean.Injection(&discoveryStatStore)
	bean.Injection(&discoveryRuntime)
	bean.Injection(&discoveryHistory)
	bean.Injection(&commonGroup)

	bean.Injection(&application)
	bean.Injection(&applicationRuntime)
	bean.Injection(&applicationVersion)
	bean.Injection(&applicationStat)
	bean.Injection(&applicationAuth)
	bean.Injection(&applicationAuthVersion)
	bean.Injection(&applicationAuthStat)
	bean.Injection(&applicationAuthRuntimeStore)
	bean.Injection(&applicationAuthPublish)
	bean.Injection(&applicationHistory)
	bean.Injection(&applicationAuthHistory)

	bean.Injection(&strategy)
	bean.Injection(&strategyStat)

	bean.Injection(&strategyVersion)

	bean.Injection(&strategyHistory)

	bean.Injection(&userInfo)
	bean.Injection(&role)
	bean.Injection(&roleAccess)
	bean.Injection(&userRole)
	bean.Injection(&roleAccessLog)

	bean.Injection(&auditLog)
	bean.Injection(&systemInfo)
	bean.Injection(&externalAPP)
	bean.Injection(&monitor)

	bean.Injection(&iNoticeChannelStore)
	bean.Injection(&iNoticeChannelStatStore)
	bean.Injection(&iNoticeChannelVersionStore)

	bean.Injection(&iWarnStrategyStore)
	bean.Injection(&iWarnHistoryStore)
}
