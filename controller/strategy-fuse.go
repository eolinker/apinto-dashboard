package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyFuseRouter(router gin.IRoutes) {
	strategyService := service.NewStrategyService(strategy_handler.NewStrategyFuseHandler("strategy-fuse"), enum.StrategyFuseRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/fuse", genAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.list)
	router.GET("/strategy/fuse", genAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.get)
	router.POST("/strategy/fuse", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyFuse), c.create)
	router.PUT("/strategy/fuse", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.update)
	router.DELETE("/strategy/fuse", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyFuse), c.del)
	router.PATCH("/strategy/fuse/restore", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.restore)
	router.PATCH("/strategy/fuse/stop", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.updateStop)
	router.GET("/strategy/fuse/to-publishs", genAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.toPublish)
	router.POST("/strategy/fuse/publish", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindStrategyFuse), c.publish)
	router.POST("/strategy/fuse/priority", genAccessHandler(access.StrategyFuseEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.changePriority)
	router.GET("/strategy/fuse/publish-history", genAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.publishHistory)
}
