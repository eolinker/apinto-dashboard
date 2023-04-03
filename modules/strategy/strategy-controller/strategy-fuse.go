package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyFuseRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyFuseHandler("strategy-fuse"), enum.StrategyFuseRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/fuse", controller.GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.list)
	router.GET("/strategy/fuse", controller.GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.get)
	router.POST("/strategy/fuse", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyFuse), c.create)
	router.PUT("/strategy/fuse", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.update)
	router.DELETE("/strategy/fuse", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyFuse), c.del)
	router.PATCH("/strategy/fuse/restore", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.restore)
	router.PATCH("/strategy/fuse/stop", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.updateStop)
	router.GET("/strategy/fuse/to-publishs", controller.GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.toPublish)
	router.POST("/strategy/fuse/publish", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyFuse), c.publish)
	router.POST("/strategy/fuse/priority", controller.GenAccessHandler(access.StrategyFuseEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.changePriority)
	router.GET("/strategy/fuse/publish-history", controller.GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.publishHistory)
}
