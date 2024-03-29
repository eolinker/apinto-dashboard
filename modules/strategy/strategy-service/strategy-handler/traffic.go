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

type trafficHandler struct {
	strategy.FormatHandler[strategy_entry.StrategyTrafficLimitConfig]
	apintoDriverName string
}

func (t *trafficHandler) GetListLabel(conf *strategy_entry.StrategyTrafficLimitConfig) string {
	limits := make([]string, 0)
	for _, v := range conf.Metrics {
		limits = append(limits, config.Metrics(v).Title())

	}
	return strings.Join(limits, ",")
}

func (t *trafficHandler) GetType() string {
	return config.StrategyTraffic
}

func (t *trafficHandler) GetConfName() string {
	return config.StrategyTrafficApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *trafficHandler) GetBatchSettingName() string {
	return config.StrategyTrafficBatchName
}

func (t *trafficHandler) CheckInput(input *strategy_dto.StrategyInfoInput[strategy_entry.StrategyTrafficLimitConfig]) error {
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

	for _, metric := range input.Config.Metrics {
		switch config.Metrics(metric) {
		case config.MetricsIP, config.MetricsAPI, config.MetricsService, config.MetricsAPP, config.MetricsStrategy:
		default:
			return fmt.Errorf("Metric %s is illegal. ", metric)
		}
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
	return checkFilters(input.Filters, "filter")
}

func (t *trafficHandler) ToApintoConfig(conf strategy_entry.StrategyTrafficLimitConfig) interface{} {
	return conf
}

func NewStrategyTrafficHandler(apintoDriverName string) strategy.IStrategyHandler[strategy_entry.StrategyTrafficLimitConfig, strategy_entry.StrategyTrafficLimitConfig] {
	return &trafficHandler{
		apintoDriverName: apintoDriverName,
	}
}
