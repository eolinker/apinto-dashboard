package strategy_handler

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"strings"
)

type cacheHandler struct {
	strategy.FormatHandler[strategy_entry.StrategyCacheConfig]
	apintoDriverName string
}

func (t *cacheHandler) GetListLabel(conf *strategy_entry.StrategyCacheConfig) string {
	return fmt.Sprintf("%v秒", conf.ValidTime)
}

func (t *cacheHandler) GetType() string {
	return enum.StrategyCache
}

func (t *cacheHandler) GetConfName() string {
	return enum.StrategyCacheApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *cacheHandler) GetBatchSettingName() string {
	return enum.StrategyCacheBatchName
}

func (t *cacheHandler) CheckInput(input *strategy_dto.StrategyInfoInput[strategy_entry.StrategyCacheConfig]) error {
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

	if input.Config.ValidTime < 1 {
		return errors.New("expire_time must be greater than 0.")
	}

	//校验筛选条件
	return checkFilters(input.Filters)
}

func (t *cacheHandler) ToApintoConfig(conf strategy_entry.StrategyCacheConfig) interface{} {
	return conf.ValidTime
}

func NewStrategyCacheHandler(apintoDriverName string) strategy.IStrategyHandler[strategy_entry.StrategyCacheConfig, strategy_entry.StrategyCacheConfig] {
	return &cacheHandler{
		apintoDriverName: apintoDriverName,
	}
}
