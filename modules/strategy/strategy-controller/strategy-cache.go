package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyCacheRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyCacheHandler("strategy-cache"), enum.StrategyCacheRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/cache", c.list)
	router.GET("/strategy/cache", c.get)
	router.POST("/strategy/cache", controller.AuditLogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyCache, c.create))
	router.PUT("/strategy/cache", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache, c.update))
	router.DELETE("/strategy/cache", controller.AuditLogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyCache, c.del))
	router.PATCH("/strategy/cache/restore", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache, c.restore))
	router.PATCH("/strategy/cache/stop", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache, c.updateStop))
	router.GET("/strategy/cache/to-publishs", c.toPublish)
	router.POST("/strategy/cache/publish", controller.AuditLogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyCache, c.publish))
	router.POST("/strategy/cache/priority", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache, c.changePriority))
	router.GET("/strategy/cache/publish-history", c.publishHistory)
}
