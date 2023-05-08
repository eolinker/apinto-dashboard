package strategy_service

import (
	"context"
	"errors"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/eosc/common/bean"
	"sort"
)

type strategyCommonService struct {
	applicationService application.IApplicationService
	apiService         api.IAPIService
	service            upstream.IService
}

func newStrategyCommonService() strategy.IStrategyCommonService {
	s := &strategyCommonService{}
	bean.Autowired(&s.applicationService)
	bean.Autowired(&s.service)
	bean.Autowired(&s.apiService)
	return s
}

func (s *strategyCommonService) GetFilterOptions(ctx context.Context, namespaceId int) ([]*strategy_model.FilterOptionsItem, error) {

	getAppKeys, err := s.applicationService.GetAppKeys(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	resList := make([]*strategy_model.FilterOptionsItem, 0)
	resList = append(resList, &strategy_model.FilterOptionsItem{
		Name:  config.FilterApplication,
		Title: "应用",
		Type:  config.FilterTypeRemote,
	}, &strategy_model.FilterOptionsItem{
		Name:  config.FilterApi,
		Title: "API",
		Type:  config.FilterTypeRemote,
	}, &strategy_model.FilterOptionsItem{
		Name:    config.FilterPath,
		Title:   "API路径",
		Type:    config.FilterTypePattern,
		Pattern: config.ApiPathRegexp,
	}, &strategy_model.FilterOptionsItem{
		Name:  config.FilterService,
		Title: "上游服务",
		Type:  config.FilterTypeRemote,
	}, &strategy_model.FilterOptionsItem{
		Name:    config.FilterIP,
		Title:   "IP",
		Type:    config.FilterTypePattern,
		Pattern: common.CIDRIpv4Exp,
	})

	method := &strategy_model.FilterOptionsItem{
		Name:  config.FilterMethod,
		Title: "API请求方式",
		Type:  config.FilterTypeStatic,
	}

	method.Options = append(method.Options, enum.HttpALL, enum.HttpGET, enum.HttpPOST, enum.HttpPUT, enum.HttpDELETE, enum.HttpPATCH, enum.HttpHEADER, enum.HttpOPTIONS)
	resList = append(resList, method)

	appKeysList := make([]*strategy_model.FilterOptionsItem, 0)
	for _, key := range getAppKeys {
		appKeysList = append(appKeysList, &strategy_model.FilterOptionsItem{
			Name:    common.SetFilterAppKey(key.Key),
			Title:   key.KeyName,
			Type:    config.FilterTypeStatic,
			Options: key.Values,
		})
	}

	sort.Slice(appKeysList, func(i, j int) bool {
		return appKeysList[i].Name > appKeysList[j].Name
	})

	resList = append(resList, appKeysList...)
	return resList, nil
}

func (s *strategyCommonService) GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*strategy_model.FilterRemoteOutput, int, error) {

	result := &strategy_model.FilterRemoteOutput{}

	switch targetType {
	case config.FilterApplication:
		applications, count, err := s.applicationService.AppListFilter(ctx, namespaceId, pageNum, pageSize, keyword)
		if err != nil {
			return nil, 0, err
		}

		for _, applicationInfo := range applications {
			result.Applications = append(result.Applications, &strategy_model.RemoteApplications{
				Name: applicationInfo.Name,
				Uuid: applicationInfo.IdStr,
				Desc: applicationInfo.Desc,
			})
		}
		result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
			Title: "应用名称",
			Field: "name",
		}, &strategy_model.RemoteTitles{
			Title: "应用ID",
			Field: "uuid",
		}, &strategy_model.RemoteTitles{
			Title: "应用描述",
			Field: "desc",
		})

		result.Target = "applications"
		return result, count, nil
	case config.FilterApi:
		remoteApis, count, err := s.apiService.GetAPIRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword, groupUUID)
		if err != nil {
			return nil, 0, err
		}

		for _, apiInfo := range remoteApis {
			result.Apis = append(result.Apis, &strategy_model.RemoteApis{
				Uuid:        apiInfo.Uuid,
				Name:        apiInfo.Name,
				Service:     apiInfo.Service,
				Group:       apiInfo.Group,
				RequestPath: apiInfo.RequestPath,
			})
		}
		result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
			Title: "API名称",
			Field: "name",
		}, &strategy_model.RemoteTitles{
			Title: "所属目录",
			Field: "group",
		}, &strategy_model.RemoteTitles{
			Title: "请求路径",
			Field: "request_path",
		})

		result.Target = "apis"
		return result, count, nil
	case config.FilterService:
		remoteServices, count, err := s.service.GetServiceRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword)
		if err != nil {
			return nil, 0, err
		}

		for _, remoteService := range remoteServices {
			result.Services = append(result.Services, &strategy_model.RemoteServices{
				Uuid:   remoteService.Uuid,
				Name:   remoteService.Name,
				Scheme: remoteService.Scheme,
				Desc:   remoteService.Desc,
			})
		}
		result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
			Title: "上游服务名称",
			Field: "name",
		}, &strategy_model.RemoteTitles{
			Title: "协议类型",
			Field: "scheme",
		}, &strategy_model.RemoteTitles{
			Title: "描述",
			Field: "desc",
		})

		result.Target = "services"
		return result, count, nil
	default:
		return nil, 0, errors.New("param error")
	}
}

func (s *strategyCommonService) GetMetricsOptions() ([]*strategy_model.MetricsOptionsItem, error) {
	resList := make([]*strategy_model.MetricsOptionsItem, 0)
	resList = append(resList, &strategy_model.MetricsOptionsItem{
		Name:  config.MetricsAPP,
		Title: "应用",
	}, &strategy_model.MetricsOptionsItem{
		Name:  config.MetricsAPI,
		Title: "API",
	}, &strategy_model.MetricsOptionsItem{
		Name:  config.MetricsService,
		Title: "上游服务",
	}, &strategy_model.MetricsOptionsItem{
		Name:  config.MetricsStrategy,
		Title: "策略",
	}, &strategy_model.MetricsOptionsItem{
		Name:  config.MetricsIP,
		Title: "IP",
	})
	return resList, nil
}
