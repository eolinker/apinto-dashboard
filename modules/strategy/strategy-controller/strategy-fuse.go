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
	router.POST("/strategy/fuse", controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyFuse), c.create)
	router.PUT("/strategy/fuse", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.update)
	router.DELETE("/strategy/fuse", controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyFuse), c.del)
	router.PATCH("/strategy/fuse/restore", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.restore)
	router.PATCH("/strategy/fuse/stop", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.updateStop)
	router.GET("/strategy/fuse/to-publishs", c.toPublish)
	router.POST("/strategy/fuse/publish", controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyFuse), c.publish)
	router.POST("/strategy/fuse/priority", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.changePriority)
	router.GET("/strategy/fuse/publish-history", c.publishHistory)
}
