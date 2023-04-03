package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyVisitRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyVisitHandler("strategy-visit"), enum.StrategyVisitRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/visit", controller.GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.list)
	router.GET("/strategy/visit", controller.GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.get)
	router.POST("/strategy/visit", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyVisit), c.create)
	router.PUT("/strategy/visit", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.update)
	router.DELETE("/strategy/visit", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyVisit), c.del)
	router.PATCH("/strategy/visit/restore", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.restore)
	router.PATCH("/strategy/visit/stop", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.updateStop)
	router.GET("/strategy/visit/to-publishs", controller.GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.toPublish)
	router.POST("/strategy/visit/publish", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyVisit), c.publish)
	router.POST("/strategy/visit/priority", controller.GenAccessHandler(access.StrategyVisitEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.changePriority)
	router.GET("/strategy/visit/publish-history", controller.GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.publishHistory)
}
