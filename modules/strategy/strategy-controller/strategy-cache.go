package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyCacheRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyCacheHandler("strategy-cache"), enum.StrategyCacheRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/cache", controller.GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.list)
	router.GET("/strategy/cache", controller.GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.get)
	router.POST("/strategy/cache", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyCache), c.create)
	router.PUT("/strategy/cache", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.update)
	router.DELETE("/strategy/cache", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyCache), c.del)
	router.PATCH("/strategy/cache/restore", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.restore)
	router.PATCH("/strategy/cache/stop", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.updateStop)
	router.GET("/strategy/cache/to-publishs", controller.GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.toPublish)
	router.POST("/strategy/cache/publish", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyCache), c.publish)
	router.POST("/strategy/cache/priority", controller.GenAccessHandler(access.StrategyCacheEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.changePriority)
	router.GET("/strategy/cache/publish-history", controller.GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.publishHistory)
}
