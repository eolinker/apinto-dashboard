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
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"sort"
)

var (
	_             apinto_module.IFilterOptionHandlerManager = (*strategyCommonService)(nil)
	staticOptions                                           = []*strategy_model.FilterOptionsItem{
		{
			Name:    config.FilterMethod,
			Title:   "API请求方式",
			Type:    config.FilterTypeStatic,
			Options: []string{enum.HttpALL, enum.HttpGET, enum.HttpPOST, enum.HttpPUT, enum.HttpDELETE, enum.HttpPATCH, enum.HttpHEADER, enum.HttpOPTIONS},
		}, {
			Name:    config.FilterPath,
			Title:   "API路径",
			Type:    config.FilterTypePattern,
			Pattern: config.ApiPathRegexp,
		}, {
			Name:    config.FilterIP,
			Title:   "IP",
			Type:    config.FilterTypePattern,
			Pattern: common.CIDRIpv4Exp,
		},
	}
)

type strategyFilterOptions struct {
	options             []*strategy_model.FilterOptionsItem
	remoteOptionHandler map[string]apinto_module.IFilterOptionHandler
}
type strategyCommonService struct {
	applicationService application.IApplicationService
	apiService         api.IAPIService
	service            upstream.IService

	filterOptions *strategyFilterOptions
}

func (s *strategyCommonService) ResetFilterOptionHandlers(handlers map[string]apinto_module.IFilterOptionHandler) {
	options := make([]*strategy_model.FilterOptionsItem, 0, len(staticOptions)+len(handlers))
	options = append(options, staticOptions...)
	if handlers != nil {
		for _, h := range handlers {
			optionConfig := h.Config()
			options = append(options, &strategy_model.FilterOptionsItem{
				Name:  optionConfig.Name,
				Title: optionConfig.Title,
				Type:  config.FilterTypeRemote,
			})

		}
	}
	sort.Sort(strategy_model.FilterOptionsItems(options))

	s.filterOptions = &strategyFilterOptions{
		options:             options,
		remoteOptionHandler: handlers,
	}

}

func newStrategyCommonService() strategy.IStrategyCommonService {
	s := &strategyCommonService{}
	bean.Autowired(&s.applicationService)
	bean.Autowired(&s.service)
	bean.Autowired(&s.apiService)
	var fm apinto_module.IFilterOptionHandlerManager = s
	bean.Injection(&fm)
	s.ResetFilterOptionHandlers(nil)
	return s
}

func (s *strategyCommonService) GetFilterOptions(ctx context.Context, namespaceId int) ([]*strategy_model.FilterOptionsItem, error) {

	return s.filterOptions.options, nil

}

func (s *strategyCommonService) GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*strategy_model.FilterRemoteOutput, int, error) {

	if s.filterOptions == nil {
		return nil, 0, errors.New("param error")
	}
	h, has := s.filterOptions.remoteOptionHandler[targetType]
	if !has {
		return nil, 0, errors.New("param error")
	}
	values, total := h.GetOptions(namespaceId, keyword, groupUUID, pageNum, pageSize)
	optionConfig := h.Config()
	return &strategy_model.FilterRemoteOutput{
		Target: "",
		Titles: optionConfig.Titles,
		Key:    optionConfig.Key,
		List:   values,
	}, total, nil
	//result := &strategy_model.FilterRemoteOutput{}
	//
	//switch targetType {
	//case config.FilterApplication:
	//	applications, count, err := s.applicationService.AppListFilter(ctx, namespaceId, pageNum, pageSize, keyword)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//
	//	for _, applicationInfo := range applications {
	//		result.Applications = append(result.Applications, &strategy_model.RemoteApplications{
	//			Name: applicationInfo.Name,
	//			Uuid: applicationInfo.IdStr,
	//			Desc: applicationInfo.Desc,
	//		})
	//	}
	//	result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
	//		Title: "应用名称",
	//		Field: "name",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "应用ID",
	//		Field: "uuid",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "应用描述",
	//		Field: "desc",
	//	})
	//
	//	result.Target = "applications"
	//	return result, count, nil
	//case config.FilterApi:
	//	remoteApis, count, err := s.apiService.GetAPIRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword, groupUUID)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//
	//	for _, apiInfo := range remoteApis {
	//		result.Apis = append(result.Apis, &strategy_model.RemoteApis{
	//			Uuid:        apiInfo.Uuid,
	//			Name:        apiInfo.Name,
	//			Service:     apiInfo.Service,
	//			Group:       apiInfo.Group,
	//			RequestPath: apiInfo.RequestPath,
	//		})
	//	}
	//	result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
	//		Title: "API名称",
	//		Field: "name",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "所属目录",
	//		Field: "group",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "请求路径",
	//		Field: "request_path",
	//	})
	//
	//	result.Target = "apis"
	//	return result, count, nil
	//case config.FilterService:
	//	remoteServices, count, err := s.service.GetServiceRemoteOptions(ctx, namespaceId, pageNum, pageSize, keyword)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//
	//	for _, remoteService := range remoteServices {
	//		result.Services = append(result.Services, &strategy_model.RemoteServices{
	//			Uuid:   remoteService.Uuid,
	//			Name:   remoteService.Name,
	//			Scheme: remoteService.Scheme,
	//			Desc:   remoteService.Desc,
	//		})
	//	}
	//	result.Titles = append(result.Titles, &strategy_model.RemoteTitles{
	//		Title: "上游服务名称",
	//		Field: "name",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "协议类型",
	//		Field: "scheme",
	//	}, &strategy_model.RemoteTitles{
	//		Title: "描述",
	//		Field: "desc",
	//	})
	//
	//	result.Target = "services"
	//	return result, count, nil
	//default:
	//	return nil, 0, errors.New("param error")
	//}
}
