package online_service

import (
	"context"
	apiService "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/application"
	clusterService "github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
	"github.com/eolinker/apinto-dashboard/modules/online"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/eosc/common/bean"
)

type resetOnlineService struct {
	list []online.IResetOnlineService
}

func (r *resetOnlineService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	for _, h := range r.list {
		h.ResetOnline(ctx, namespaceId, clusterId)
	}
}

func newResetOnline() online.IResetOnlineService {
	onlineService := &resetOnlineService{}
	var clConfig clusterService.IClusterConfigService
	var variableService variable.IClusterVariableService
	var discoveryService discovery.IDiscoveryService
	var iService upstream.IService
	var api apiService.IAPIService
	var applicationService application.IApplicationService
	var commonStrategy strategy.IStrategyCommonService

	bean.Autowired(&clConfig)
	bean.Autowired(&variableService)
	bean.Autowired(&api)
	bean.Autowired(&iService)
	bean.Autowired(&discoveryService)
	bean.Autowired(&applicationService)
	//bean.Autowired(&strategy)
	bean.Autowired(&commonStrategy)
	onlineService.list = append(onlineService.list, clConfig, variableService, discoveryService, iService, api, applicationService, commonStrategy)
	return onlineService
}
