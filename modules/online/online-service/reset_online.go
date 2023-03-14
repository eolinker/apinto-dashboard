package online_service

import (
	"context"
	api2 "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/application"
	cluster_service "github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
	"github.com/eolinker/apinto-dashboard/modules/online"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
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
	online := &resetOnlineService{}
	var clConfig cluster_service.IClusterConfigService
	var variable variable.IClusterVariableService
	var discovery discovery.IDiscoveryService
	var iService upstream.IService
	var api api2.IAPIService
	var application application.IApplicationService
	var commonStrategy strategy.IStrategyCommonService

	bean.Autowired(&clConfig)
	bean.Autowired(&variable)
	bean.Autowired(&api)
	bean.Autowired(&iService)
	bean.Autowired(&discovery)
	bean.Autowired(&application)
	//bean.Autowired(&strategy)
	bean.Autowired(&commonStrategy)
	online.list = append(online.list, clConfig, variable, discovery, iService, api, application, commonStrategy)
	return online
}
