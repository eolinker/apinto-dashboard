package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyTrafficRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyTrafficHandler("strategy-limiting"), enum.StrategyTrafficRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/traffic", c.list)
	router.GET("/strategy/traffic", c.get)
	router.POST("/strategy/traffic", controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyTraffic), c.create)
	router.PUT("/strategy/traffic", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.update)
	router.DELETE("/strategy/traffic", controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyTraffic), c.del)
	router.PATCH("/strategy/traffic/restore", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.restore)
	router.PATCH("/strategy/traffic/stop", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.updateStop)
	router.GET("/strategy/traffic/to-publishs", c.toPublish)
	router.POST("/strategy/traffic/publish", controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyTraffic), c.publish)
	router.POST("/strategy/traffic/priority", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.changePriority)
	router.GET("/strategy/traffic/publish-history", c.publishHistory)
}
