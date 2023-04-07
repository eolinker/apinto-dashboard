package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyVisitRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyVisitHandler("strategy-visit"), enum.StrategyVisitRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/visit", c.list)
	router.GET("/strategy/visit", c.get)
	router.POST("/strategy/visit", controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyVisit), c.create)
	router.PUT("/strategy/visit", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.update)
	router.DELETE("/strategy/visit", controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyVisit), c.del)
	router.PATCH("/strategy/visit/restore", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.restore)
	router.PATCH("/strategy/visit/stop", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.updateStop)
	router.GET("/strategy/visit/to-publishs", c.toPublish)
	router.POST("/strategy/visit/publish", controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyVisit), c.publish)
	router.POST("/strategy/visit/priority", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.changePriority)
	router.GET("/strategy/visit/publish-history", c.publishHistory)
}
