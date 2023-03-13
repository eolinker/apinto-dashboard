package service

import (
	"context"
	"errors"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/eosc/common/bean"
	"sort"
)

type IStrategyCommonService interface {
	GetFilterOptions(ctx context.Context, namespaceId int) ([]*model.FilterOptionsItem, error)
	GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*model.FilterRemoteOutput, int, error)
	GetMetricsOptions() ([]*model.MetricsOptionsItem, error)
	AddHandler(onlineService IResetOnlineService)
	IResetOnlineService
}

type strategyCommonService struct {
	applicationService  IApplicationService
	apiService          IAPIService
	service             IService
	resetOnlineServices []IResetOnlineService
}

func (s *strategyCommonService) AddHandler(onlineService IResetOnlineService) {
	s.resetOnlineServices = append(s.resetOnlineServices, onlineService)
}

func newStrategyCommonService() IStrategyCommonService {
	s := &strategyCommonService{}
	bean.Autowired(&s.applicationService)
	bean.Autowired(&s.service)
	bean.Autowired(&s.apiService)
	return s
}

func (s *strategyCommonService) GetFilterOptions(ctx context.Context, namespaceId int) ([]*model.FilterOptionsItem, error) {

	getAppKeys, err := s.applicationService.GetAppKeys(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	resList := make([]*model.FilterOptionsItem, 0)
	resList = append(resList, &model.FilterOptionsItem{
		Name:  enum.FilterApplication,
		Title: "应用",
		Type:  enum.FilterTypeRemote,
	}, &model.FilterOptionsItem{
		Name:  enum.FilterApi,
		Title: "API",
		Type:  enum.FilterTypeRemote,
	}, &model.FilterOptionsItem{
		Name:    enum.FilterPath,
		Title:   "API路径",
		Type:    enum.FilterTypePattern,
		Pattern: enum.ApiPathRegexp,
	}, &model.FilterOptionsItem{
		Name:  enum.FilterService,
		Title: "上游服务",
		Type:  enum.FilterTypeRemote,
	}, &model.FilterOptionsItem{
		Name:    enum.FilterIP,
		Title:   "IP",
		Type:    enum.FilterTypePattern,
		Pattern: common.CIDRIpv4Exp,
	})

	method := &model.FilterOptionsItem{
		Name:  enum.FilterMethod,
		Title: "API请求方式",
		Type:  enum.FilterTypeStatic,
	}

	method.Options = append(method.Options, enum.HttpALL, enum.HttpGET, enum.HttpPOST, enum.HttpPUT, enum.HttpDELETE, enum.HttpPATCH, enum.HttpHEADER, enum.HttpOPTIONS)
	resList = append(resList, method)

	appKeysList := make([]*model.FilterOptionsItem, 0)
	for _, key := range getAppKeys {
		appKeysList = append(appKeysList, &model.FilterOptionsItem{
			Name:    common.SetFilterAppKey(key.Key),
			Title:   key.KeyName,
			Type:    enum.FilterTypeStatic,
			Options: key.Values,
		})
	}

	sort.Slice(appKeysList, func(i, j int) bool {
		return appKeysList[i].Name > appKeysList[j].Name
	})

	resList = append(resList, appKeysList...)
	return resList, nil
}

func (s *strategyCommonService) GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*model.FilterRemoteOutput, int, error) {

	result := &model.FilterRemoteOutput{}

	switch targetType {
	case enum.FilterApplication:
		applications, count, err := s.applicationService.AppListFilter(ctx, namespaceId, pageNum, pageSize, keyword)
		if err != nil {
			return nil, 0, err
		}

		for _, application := range applications {
			result.Applications = append(result.Applications, &model.RemoteApplications{
				Name: application.Name,
				Uuid: application.IdStr,
				Desc: application.Desc,
			})
		}
		result.Titles = append(result.Titles, &model.RemoteTitles{
			Title: "应用名称",
			Field: "name",
		}, &model.RemoteTitles{
			Title: "应用ID",
			Field: "uuid",
		}, &model.RemoteTitles{
			Title: "应用描述",
			Field: "desc",
		})

		result.Target = "applications"
		return result, count, nil
	case enum.FilterApi:
		remoteApis, count, err := s.apiService.GetAPIRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword, groupUUID)
		if err != nil {
			return nil, 0, err
		}

		for _, api := range remoteApis {
			result.Apis = append(result.Apis, &model.RemoteApis{
				Uuid:        api.Uuid,
				Name:        api.Name,
				Service:     api.Service,
				Group:       api.Group,
				RequestPath: api.RequestPath,
			})
		}
		result.Titles = append(result.Titles, &model.RemoteTitles{
			Title: "API名称",
			Field: "name",
		}, &model.RemoteTitles{
			Title: "所属目录",
			Field: "group",
		}, &model.RemoteTitles{
			Title: "请求路径",
			Field: "request_path",
		})

		result.Target = "apis"
		return result, count, nil
	case enum.FilterService:
		remoteServices, count, err := s.service.GetServiceRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword)
		if err != nil {
			return nil, 0, err
		}

		for _, remoteService := range remoteServices {
			result.Services = append(result.Services, &model.RemoteServices{
				Uuid:   remoteService.Uuid,
				Name:   remoteService.Name,
				Scheme: remoteService.Scheme,
				Desc:   remoteService.Desc,
			})
		}
		result.Titles = append(result.Titles, &model.RemoteTitles{
			Title: "上游服务名称",
			Field: "name",
		}, &model.RemoteTitles{
			Title: "协议类型",
			Field: "scheme",
		}, &model.RemoteTitles{
			Title: "描述",
			Field: "desc",
		})

		result.Target = "services"
		return result, count, nil
	default:
		return nil, 0, errors.New("param error")
	}
}

func (s *strategyCommonService) GetMetricsOptions() ([]*model.MetricsOptionsItem, error) {
	resList := make([]*model.MetricsOptionsItem, 0)
	resList = append(resList, &model.MetricsOptionsItem{
		Name:  enum.MetricsAPP,
		Title: "应用",
	}, &model.MetricsOptionsItem{
		Name:  enum.MetricsAPI,
		Title: "API",
	}, &model.MetricsOptionsItem{
		Name:  enum.MetricsService,
		Title: "上游服务",
	}, &model.MetricsOptionsItem{
		Name:  enum.MetricsStrategy,
		Title: "策略",
	}, &model.MetricsOptionsItem{
		Name:  enum.MetricsIP,
		Title: "IP",
	})
	return resList, nil
}

func (s *strategyCommonService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	for _, onlineService := range s.resetOnlineServices {
		onlineService.ResetOnline(ctx, namespaceId, clusterId)
	}
}
