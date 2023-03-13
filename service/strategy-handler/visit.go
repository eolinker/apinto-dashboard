package strategy_handler

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"strings"
)

type visitHandler struct {
	applicationService service.IApplicationService
	apiService         service.IAPIService
	service            service.IService
	apintoDriverName   string
}

func (t *visitHandler) GetListLabel(conf *entry.StrategyVisitConfig) string {
	switch conf.VisitRule {
	case enum.VisitRuleAllow:
		return "允许访问"
	case enum.VisitRuleRefuse:
		return "拒绝访问"
	default:
		return ""
	}
}

func (t *visitHandler) GetType() string {
	return enum.StrategyVisit
}

func (t *visitHandler) GetConfName() string {
	return enum.StrategyVisitApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *visitHandler) GetBatchSettingName() string {
	return enum.StrategyVisitBatchName
}

func (t *visitHandler) CheckInput(input *dto.StrategyInfoInput[entry.StrategyVisitConfig]) error {
	input.Uuid = strings.TrimSpace(input.Uuid)
	if input.Uuid != "" {
		err := common.IsMatchString(common.UUIDExp, input.Uuid)
		if err != nil {
			return err
		}
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return errors.New("Name can't be null. ")
	}
	if input.Priority < 0 {
		input.Priority = 0
	}

	if input.Config == nil {
		return errors.New("config can't be null. ")
	}

	//检查生效规则
	switch input.Config.VisitRule {
	case enum.VisitRuleAllow, enum.VisitRuleRefuse:
	default:
		return fmt.Errorf("visit_rule %s is illegal. ", input.Config.VisitRule)
	}

	//检查生效范围
	filterNameSet := make(map[string]struct{})
	for _, influence := range input.Config.InfluenceSphere {
		switch influence.Name {
		case enum.FilterApplication, enum.FilterApi, enum.FilterPath, enum.FilterService, enum.FilterMethod, enum.FilterIP:
		default:
			if !common.IsMatchFilterAppKey(influence.Name) {
				return fmt.Errorf("influence_sphere.Name %s is illegal. ", influence.Name)
			}
		}

		if len(influence.Values) == 0 {
			return fmt.Errorf("influence_sphere.Options can't be null. influence_sphere.Name:%s ", influence.Name)
		}

		if _, has := filterNameSet[influence.Name]; has {
			return fmt.Errorf("influenceName %s is reduplicative. ", influence.Name)
		}
		filterNameSet[influence.Name] = struct{}{}
	}

	//校验筛选条件
	return checkFilters(input.Filters)
}

func (t *visitHandler) ToApintoConfig(conf entry.StrategyVisitConfig) interface{} {
	influenceSphere := make(map[string][]string)

	for _, filter := range conf.InfluenceSphere {
		influenceSphere[filter.Name] = filter.Values
	}

	return &v1.StrategyVisit{
		VisitRule:       conf.VisitRule,
		InfluenceSphere: influenceSphere,
		Continue:        conf.Continue,
	}
}

func (t *visitHandler) FormatOut(ctx context.Context, namespaceID int, input *model.StrategyInfoOutput[entry.StrategyVisitConfig]) *model.StrategyInfoOutput[model.VisitInfoOutputConf] {
	output := new(model.StrategyInfoOutput[model.VisitInfoOutputConf])
	output.Strategy = input.Strategy
	output.Filters = input.Filters

	config := &model.VisitInfoOutputConf{
		VisitRule:       input.Config.VisitRule,
		InfluenceSphere: nil,
		Continue:        input.Config.Continue,
		Extender:        nil,
	}

	//拼装生效范围以及extender
	extenderData := &model.ExtenderData{
		Api:         make(map[string]*model.RemoteApis),
		Service:     make(map[string]*model.RemoteServices),
		Application: make(map[string]*model.RemoteApplications),
	}
	filters := make([]*model.FilterOutput, 0, len(input.Filters))
	for _, f := range input.Config.InfluenceSphere {
		filter := &model.FilterOutput{
			Name:   f.Name,
			Values: f.Values,
		}
		if len(f.Values) == 0 {
			continue
		}
		switch filter.Name {
		case enum.FilterApplication:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "应用"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有应用"
			} else {
				apps, _ := t.applicationService.AppListByUUIDS(ctx, namespaceID, f.Values)
				if len(apps) == 0 {
					continue
				}
				labels := make([]string, len(apps))
				for i, app := range apps {
					extenderData.Application[app.IdStr] = &model.RemoteApplications{
						Name: app.Name,
						Uuid: app.IdStr,
						Desc: app.Desc,
					}
					labels[i] = app.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterApi:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "API"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有API"
			} else {
				apis, _ := t.apiService.GetAPIRemoteByUUIDS(ctx, namespaceID, f.Values)
				if len(apis) == 0 {
					continue
				}
				labels := make([]string, len(apis))
				for i, api := range apis {
					extenderData.Api[api.Uuid] = api
					labels[i] = api.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterPath:
			filter.Type = enum.FilterTypePattern
			filter.Title = "API路径"
			filter.Label = f.Values[0]
		case enum.FilterService:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "上游服务"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有上游服务"
			} else {
				services, _ := t.service.GetServiceRemoteByNames(ctx, namespaceID, f.Values)
				if len(services) == 0 {
					continue
				}
				labels := make([]string, len(services))
				for i, sv := range services {
					extenderData.Service[sv.Uuid] = sv
					labels[i] = sv.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterMethod:
			filter.Type = enum.FilterTypeStatic
			filter.Title = "API请求方式"
			if f.Values[0] == enum.HttpALL {
				filter.Label = "所有API请求方式"
			} else {
				filter.Label = strings.Join(f.Values, ",")
			}
		case enum.FilterIP:
			filter.Type = enum.FilterTypePattern
			filter.Title = "IP"
			filter.Label = strings.Join(f.Values, ",")
		default: //KEY(应用)
			filter.Type = enum.FilterTypeStatic
			filter.Title = fmt.Sprintf("%s(应用)", common.GetFilterAppKey(filter.Name))
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = fmt.Sprintf("%s(应用)所有值", common.GetFilterAppKey(filter.Name))
			} else {
				filter.Label = strings.Join(f.Values, ",")
			}
		}
		filters = append(filters, filter)
	}

	config.InfluenceSphere = filters
	config.Extender = extenderData

	output.Config = config
	return output
}

func NewStrategyVisitHandler(apintoDriverName string) service.IStrategyHandler[entry.StrategyVisitConfig, model.VisitInfoOutputConf] {
	h := &visitHandler{
		apintoDriverName: apintoDriverName,
	}
	bean.Autowired(&h.service)
	bean.Autowired(&h.apiService)
	bean.Autowired(&h.applicationService)

	return h
}
