package strategy_handler

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"strings"
)

type cacheHandler struct {
	service.FormatHandler[entry.StrategyCacheConfig]
	apintoDriverName string
}

func (t *cacheHandler) GetListLabel(conf *entry.StrategyCacheConfig) string {
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

func (t *cacheHandler) CheckInput(input *dto.StrategyInfoInput[entry.StrategyCacheConfig]) error {
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

func (t *cacheHandler) ToApintoConfig(conf entry.StrategyCacheConfig) interface{} {
	return conf.ValidTime
}

func NewStrategyCacheHandler(apintoDriverName string) service.IStrategyHandler[entry.StrategyCacheConfig, entry.StrategyCacheConfig] {
	return &cacheHandler{
		apintoDriverName: apintoDriverName,
	}
}
