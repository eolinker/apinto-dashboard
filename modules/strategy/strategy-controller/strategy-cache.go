package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyCacheController() *strategyController[strategy_entry.StrategyCacheConfig, strategy_entry.StrategyCacheConfig] {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyCacheHandler("strategy-cache"), config.StrategyCacheRuntimeKind)

	c := newStrategyController(strategyService)
	return c
}
