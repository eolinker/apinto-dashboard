package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyTrafficController() *strategyController[strategy_entry.StrategyTrafficLimitConfig, strategy_entry.StrategyTrafficLimitConfig] {
	handler := strategy_handler.NewStrategyTrafficHandler("strategy-limiting")
	strategyService := strategy_service.NewStrategyService(handler, config.StrategyTrafficRuntimeKind)

	c := newStrategyController(strategyService, handler.GetType())

	return c
}
