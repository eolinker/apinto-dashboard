package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyGreyController() *strategyController[strategy_entry.StrategyGreyConfig, strategy_entry.StrategyGreyConfig] {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyGreyHandler("strategy-grey"), config.StrategyGreyRuntimeKind)

	c := newStrategyController(strategyService)
	return c
}
