package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/enum"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyTrafficController() *strategyController[strategy_entry.StrategyTrafficLimitConfig, strategy_entry.StrategyTrafficLimitConfig] {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyTrafficHandler("strategy-limiting"), enum.StrategyTrafficRuntimeKind)

	c := newStrategyController(strategyService)

	return c
}
