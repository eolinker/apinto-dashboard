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
	router.GET("/strategies/fuse", GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.list)
	router.GET("/strategy/fuse", GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.get)
	router.POST("/strategy/fuse", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyFuse), c.create)
	router.PUT("/strategy/fuse", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.update)
	router.DELETE("/strategy/fuse", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyFuse), c.del)
	router.PATCH("/strategy/fuse/restore", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.restore)
	router.PATCH("/strategy/fuse/stop", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.updateStop)
	router.GET("/strategy/fuse/to-publishs", GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.toPublish)
	router.POST("/strategy/fuse/publish", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyFuse), c.publish)
	router.POST("/strategy/fuse/priority", GenAccessHandler(access.StrategyFuseEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyFuse), c.changePriority)
	router.GET("/strategy/fuse/publish-history", GenAccessHandler(access.StrategyFuseView, access.StrategyFuseEdit), c.publishHistory)
}
