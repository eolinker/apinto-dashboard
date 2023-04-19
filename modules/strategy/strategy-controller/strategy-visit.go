package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/enum"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	strategy_model "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyVisitController() *strategyController[strategy_entry.StrategyVisitConfig, strategy_model.VisitInfoOutputConf] {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyVisitHandler("strategy-visit"), enum.StrategyVisitRuntimeKind)

	c := newStrategyController(strategyService)
	return c
}
