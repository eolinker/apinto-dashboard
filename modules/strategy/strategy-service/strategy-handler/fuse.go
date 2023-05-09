package strategy_handler

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"strings"
)

type fuseHandler struct {
	strategy.FormatHandler[strategy_entry.StrategyFuseConfig]
	apintoDriverName string
}

func (t *fuseHandler) GetListLabel(conf *strategy_entry.StrategyFuseConfig) string {
	return config.Metrics(conf.Metric).Title()

}

func (t *fuseHandler) GetType() string {
	return config.StrategyFuse
}

func (t *fuseHandler) GetConfName() string {
	return config.StrategyFuseApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *fuseHandler) GetBatchSettingName() string {
	return config.StrategyFuseBatchName
}

func (t *fuseHandler) CheckInput(input *strategy_dto.StrategyInfoInput[strategy_entry.StrategyFuseConfig]) error {
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

	//校验熔断维度
	switch config.Metrics(input.Config.Metric) {
	case config.MetricsAPI, config.MetricsService:
	default:
		return fmt.Errorf("Metric %s is illegal. ", input.Config.Metric)
	}

	//校验熔断条件
	if len(input.Config.FuseCondition.StatusCodes) == 0 {
		return errors.New("fuse_condition.status_codes can't be null. ")
	}
	if err := checkStatusCode(input.Config.FuseCondition.StatusCodes...); err != nil {
		return err
	}
	if input.Config.FuseCondition.Count <= 0 {
		return errors.New("fuse_condition.count must larger than 0. ")
	}

	//校验熔断时间
	if input.Config.FuseTime.Time < 1 {
		return errors.New("fuse_time.time must larger than 0. ")
	}
	if input.Config.FuseTime.MaxTime < 1 {
		return errors.New("fuse_time.max_time must larger than 0. ")
	}

	//校验恢复条件
	if len(input.Config.RecoverCondition.StatusCodes) == 0 {
		return errors.New("recover_condition.status_codes can't be null. ")
	}
	if err := checkStatusCode(input.Config.RecoverCondition.StatusCodes...); err != nil {
		return err
	}
	if input.Config.FuseCondition.Count <= 0 {
		return errors.New("recover_condition.count must larger than 0. ")
	}

	//校验响应内容
	if err := checkStatusCode(input.Config.Response.StatusCode); err != nil {
		return err
	}
	if input.Config.Response.ContentType == "" {
		return errors.New("response.content_type can't be null. ")
	}
	//check Charset
	if err := checkCharset(input.Config.Response.Charset); err != nil {
		return err
	}

	//校验筛选条件
	return checkFilters(input.Filters)
}

func (t *fuseHandler) ToApintoConfig(conf strategy_entry.StrategyFuseConfig) interface{} {
	return conf
}

func NewStrategyFuseHandler(apintoDriverName string) strategy.IStrategyHandler[strategy_entry.StrategyFuseConfig, strategy_entry.StrategyFuseConfig] {
	return &fuseHandler{
		apintoDriverName: apintoDriverName,
	}
}
