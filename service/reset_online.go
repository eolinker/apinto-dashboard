package service

import (
	"context"
	"github.com/eolinker/eosc/common/bean"
)

type IResetOnlineService interface {
	ResetOnline(ctx context.Context, namespaceId, clusterId int)
}

type resetOnlineService struct {
	list []IResetOnlineService
}

func (r *resetOnlineService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	for _, online := range r.list {
		online.ResetOnline(ctx, namespaceId, clusterId)
	}
}

func newResetOnline() IResetOnlineService {
	online := &resetOnlineService{}
	var clConfig IClusterConfigService
	var variable IClusterVariableService
	var discovery IDiscoveryService
	var iService IService
	var api IAPIService
	var application IApplicationService
	var commonStrategy IStrategyCommonService

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
