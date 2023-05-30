package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyFuseController() *strategyController[strategy_entry.StrategyFuseConfig, strategy_entry.StrategyFuseConfig] {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyFuseHandler("strategy-fuse"), config.StrategyFuseRuntimeKind)

	c := newStrategyController(strategyService)
	return c
}
