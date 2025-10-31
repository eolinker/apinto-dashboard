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
	"strconv"
	"strings"
)

type dataMaskHandler struct {
	commonService strategy.IStrategyCommonService
	driver        string
}

func (t *dataMaskHandler) GetListLabel(conf *strategy_entry.DataMaskConfig) string {
	return strconv.Itoa(len(conf.Rules))
}

func (t *dataMaskHandler) GetType() string {
	return strategyConfig.StrategyDataMask
}

func (t *dataMaskHandler) GetConfName() string {
	return strategyConfig.StrategyDataMaskApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *dataMaskHandler) GetBatchSettingName() string {
	return strategyConfig.StrategyDataMaskBatchName
}

var validMatchInnerValues = map[string]struct{}{
	"name":      {},
	"phone":     {},
	"email":     {},
	"id-card":   {},
	"bank-card": {},
	"date":      {},
	"amount":    {},
}
var validMatchTypes = map[string]struct{}{
	"inner":     {},
	"keyword":   {},
	"regex":     {},
	"json_path": {},
}

var validMaskTypes = map[string]struct{}{
	"partial-display": {},
	"partial-masking": {},
	"truncation":      {},
	"replacement":     {},
	"shuffling":       {},
}

var validReplaceTypes = map[string]struct{}{
	"random": {},
	"custom": {},
}

func (t *dataMaskHandler) CheckInput(input *strategy_dto.StrategyInfoInput[strategy_entry.DataMaskConfig]) error {
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

	if input.Config == nil || len(input.Config.Rules) == 0 {
		return errors.New("strategyConfig can't be null. ")
	}

	for _, r := range input.Config.Rules {
		if r.Match == nil {
			return errors.New("match can't be null. ")
		}
		if r.Mask == nil {
			return errors.New("mask can't be null. ")
		}

		if _, ok := validMatchTypes[r.Match.Type]; !ok {
			return fmt.Errorf("match type %s is illegal. ", r.Match.Type)
		}

		if r.Match.Type == "inner" {
			if _, ok := validMatchInnerValues[r.Match.Value]; !ok {
				return fmt.Errorf("match value %s is illegal. ", r.Match.Value)
			}
		}
		if _, ok := validMaskTypes[r.Mask.Type]; !ok {
			return fmt.Errorf("mask type %s is illegal. ", r.Mask.Type)
		}
		if r.Mask.Replace != nil {
			if _, ok := validReplaceTypes[r.Mask.Replace.Type]; !ok {
				return fmt.Errorf("replace type %s is illegal. ", r.Mask.Replace.Type)
			}
		}
	}

	//校验筛选条件
	return checkFilters(input.Filters, "filter")
}

func (t *dataMaskHandler) ToApintoConfig(conf strategy_entry.DataMaskConfig) interface{} {
	rules := make([]*v1.DataMaskRule, 0, len(conf.Rules))
	for _, r := range conf.Rules {
		var match *v1.DataMaskBasicItem
		if r.Match != nil {
			match = &v1.DataMaskBasicItem{
				Type:  r.Match.Type,
				Value: r.Match.Value,
			}
		}
		var mask *v1.DataMaskMask
		if r.Mask != nil {
			mask = &v1.DataMaskMask{
				Type:   r.Mask.Type,
				Begin:  r.Mask.Begin,
				Length: r.Mask.Length,
			}
			if r.Mask.Replace != nil {
				mask.Replace = &v1.DataMaskBasicItem{
					Type:  r.Mask.Replace.Type,
					Value: r.Mask.Replace.Value,
				}
			}
		}
		rules = append(rules, &v1.DataMaskRule{
			Match: match,
			Mask:  mask,
		})
	}
	return &v1.StrategyDataMask{
		Rules: rules,
	}
}

func (t *dataMaskHandler) FormatOut(ctx context.Context, namespaceID int, input *strategy_model.StrategyInfoOutput[strategy_entry.DataMaskConfig]) *strategy_model.StrategyInfoOutput[strategy_model.DataMaskConfig] {
	output := new(strategy_model.StrategyInfoOutput[strategy_model.DataMaskConfig])
	output.Strategy = input.Strategy
	output.Filters = input.Filters
	rules := make([]*strategy_model.DataMaskRule, 0, len(input.Config.Rules))
	for _, r := range input.Config.Rules {
		var match *strategy_model.DataMaskBasicItem
		if r.Match != nil {
			match = &strategy_model.DataMaskBasicItem{
				Type:  r.Match.Type,
				Value: r.Match.Value,
			}
		}
		var mask *strategy_model.DataMaskMask
		if r.Mask != nil {
			mask = &strategy_model.DataMaskMask{
				Type:   r.Mask.Type,
				Begin:  r.Mask.Begin,
				Length: r.Mask.Length,
			}
			if r.Mask.Replace != nil {
				mask.Replace = &strategy_model.DataMaskBasicItem{
					Type:  r.Mask.Replace.Type,
					Value: r.Mask.Replace.Value,
				}
			}
		}
		rules = append(rules, &strategy_model.DataMaskRule{
			Match: match,
			Mask:  mask,
		})
	}
	config := &strategy_model.DataMaskConfig{
		Rules: rules,
	}

	output.Config = config
	return output
}

func NewStrategyDataMaskHandler(driver string) strategy.IStrategyHandler[strategy_entry.DataMaskConfig, strategy_model.DataMaskConfig] {
	h := &dataMaskHandler{
		driver: driver,
	}
	bean.Autowired(&h.commonService)

	return h
}
