package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyFuseRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyFuseHandler("strategy-fuse"), enum.StrategyFuseRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/fuse", c.list)
	router.GET("/strategy/fuse", c.get)
	router.POST("/strategy/fuse", controller.AuditLogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyFuse, c.create))
	router.PUT("/strategy/fuse", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse, c.update))
	router.DELETE("/strategy/fuse", controller.AuditLogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyFuse, c.del))
	router.PATCH("/strategy/fuse/restore", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse, c.restore))
	router.PATCH("/strategy/fuse/stop", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse, c.updateStop))
	router.GET("/strategy/fuse/to-publishs", c.toPublish)
	router.POST("/strategy/fuse/publish", controller.AuditLogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyFuse, c.publish))
	router.POST("/strategy/fuse/priority", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse, c.changePriority))
	router.GET("/strategy/fuse/publish-history", c.publishHistory)
}
