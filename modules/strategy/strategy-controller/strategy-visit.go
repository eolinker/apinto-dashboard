package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	strategy_entry "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	strategy_model "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
)

func newStrategyVisitController() *strategyController[strategy_entry.StrategyVisitConfig, strategy_model.VisitInfoOutputConf] {
	handler := strategy_handler.NewStrategyVisitHandler("strategy-visit")
	strategyService := strategy_service.NewStrategyService(handler, config.StrategyVisitRuntimeKind)

	c := newStrategyController(strategyService, handler.GetType())
	return c
}
