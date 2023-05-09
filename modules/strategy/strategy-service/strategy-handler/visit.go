package strategy_handler

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	strategyConfig "github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/eosc/common/bean"
	"strings"
)

type visitHandler struct {
	commonService    strategy.IStrategyCommonService
	apintoDriverName string
}

func (t *visitHandler) GetListLabel(conf *strategy_entry.StrategyVisitConfig) string {
	switch conf.VisitRule {
	case strategyConfig.VisitRuleAllow:
		return "允许访问"
	case strategyConfig.VisitRuleRefuse:
		return "拒绝访问"
	default:
		return ""
	}
}

func (t *visitHandler) GetType() string {
	return strategyConfig.StrategyVisit
}

func (t *visitHandler) GetConfName() string {
	return strategyConfig.StrategyVisitApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *visitHandler) GetBatchSettingName() string {
	return strategyConfig.StrategyVisitBatchName
}

func (t *visitHandler) CheckInput(input *strategy_dto.StrategyInfoInput[strategy_entry.StrategyVisitConfig]) error {
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
		return errors.New("strategyConfig can't be null. ")
	}

	//检查生效规则
	switch input.Config.VisitRule {
	case strategyConfig.VisitRuleAllow, strategyConfig.VisitRuleRefuse:
	default:
		return fmt.Errorf("visit_rule %s is illegal. ", input.Config.VisitRule)
	}

	//检查生效范围
	filterNameSet := make(map[string]struct{})
	for _, influence := range input.Config.InfluenceSphere {
		switch influence.Name {
		case strategyConfig.FilterApplication, strategyConfig.FilterApi, strategyConfig.FilterPath, strategyConfig.FilterService, strategyConfig.FilterMethod, strategyConfig.FilterIP:
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

func (t *visitHandler) ToApintoConfig(conf strategy_entry.StrategyVisitConfig) interface{} {
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

func (t *visitHandler) FormatOut(ctx context.Context, namespaceID int, input *strategy_model.StrategyInfoOutput[strategy_entry.StrategyVisitConfig]) *strategy_model.StrategyInfoOutput[strategy_model.VisitInfoOutputConf] {
	output := new(strategy_model.StrategyInfoOutput[strategy_model.VisitInfoOutputConf])
	output.Strategy = input.Strategy
	output.Filters = input.Filters

	config := &strategy_model.VisitInfoOutputConf{
		VisitRule:       input.Config.VisitRule,
		InfluenceSphere: nil,
		Continue:        input.Config.Continue,
	}

	filters := make([]*strategy_model.FilterOutput, 0, len(input.Filters))
	for _, f := range input.Config.InfluenceSphere {
		filter := &strategy_model.FilterOutput{
			Name:   f.Name,
			Values: f.Values,
		}
		if len(f.Values) == 0 {
			continue
		}
		filter.Title, filter.Label, filter.Type = t.commonService.GetFilterLabel(ctx, namespaceID, filter.Name, filter.Values)
		if filter.Label == "" {
			continue
		}
		filters = append(filters, filter)
	}

	config.InfluenceSphere = filters

	output.Config = config
	return output
}

func NewStrategyVisitHandler(apintoDriverName string) strategy.IStrategyHandler[strategy_entry.StrategyVisitConfig, strategy_model.VisitInfoOutputConf] {
	h := &visitHandler{
		apintoDriverName: apintoDriverName,
	}
	bean.Autowired(&h.commonService)

	return h
}
