package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyGreyController() *strategyController[strategy_entry.StrategyGreyConfig, strategy_entry.StrategyGreyConfig] {
	handler := strategy_handler.NewStrategyGreyHandler("strategy-grey")
	strategyService := strategy_service.NewStrategyService(handler, config.StrategyGreyRuntimeKind)

	c := newStrategyController(strategyService, handler.GetType())
	return c
}
