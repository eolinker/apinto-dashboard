package strategy_service

import (
	"context"
	"errors"
	"fmt"
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
	"strings"
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
	titleConfigs        map[string]string
	typeConfigs         map[string]string
}
type strategyCommonService struct {
	applicationService application.IApplicationService
	apiService         api.IAPIService
	service            upstream.IService

	filterOptions *strategyFilterOptions
}

func (s *strategyCommonService) GetFilterLabel(ctx context.Context, namespaceId int, name string, values []string) (string, string, string) {
	title, has := s.filterOptions.titleConfigs[name]
	if has {
		if len(values) > 0 {
			if values[0] == config.FilterValuesALL {
				return title, fmt.Sprintf("全部%s", title), ""
			}

			handler := s.filterOptions.remoteOptionHandler[name]
			return title, strings.Join(handler.Labels(namespaceId, values...), ","), ""
		}
		return "", "", ""
	}

	switch name {

	case config.FilterMethod:
		if len(values) > 0 {
			if values[0] == config.FilterValuesALL {
				return "API请求方式", "全部请求方式", ""

			} else {
				return "API请求方式", strings.Join(values, ","), ""

			}
		}
	case config.FilterPath:
		if len(values) > 0 {
			return "API路径", values[0], ""

		}
	case config.FilterIP:
		if len(values) > 0 {
			return "IP", strings.Join(values, ","), ""
		}
	}
	return "", "", ""
}

func (s *strategyCommonService) ResetFilterOptionHandlers(handlers map[string]apinto_module.IFilterOptionHandler) {
	options := make([]*strategy_model.FilterOptionsItem, 0, len(staticOptions)+len(handlers))
	configs := map[string]string{}
	options = append(options, staticOptions...)

	if handlers != nil {
		for name, h := range handlers {
			optionConfig := h.Config()
			options = append(options, &strategy_model.FilterOptionsItem{
				Name:  name,
				Title: optionConfig.Title,
				Type:  config.FilterTypeRemote,
			})
			configs[name] = optionConfig.Title
		}
	}
	sort.Sort(strategy_model.FilterOptionsItems(options))

	s.filterOptions = &strategyFilterOptions{
		options:             options,
		remoteOptionHandler: handlers,
		titleConfigs:        configs,
		typeConfigs: common.SliceToMapO(options, func(t *strategy_model.FilterOptionsItem) (string, string) {
			return t.Name, t.Type
		}),
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
		Target: "list",
		Titles: optionConfig.Titles,
		Key:    optionConfig.Key,
		List:   values,
	}, total, nil

}
