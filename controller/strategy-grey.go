package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyGreyRouter(router gin.IRoutes) {
	strategyService := service.NewStrategyService(strategy_handler.NewStrategyGreyHandler("strategy-grey"), enum.StrategyGreyRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/grey", GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.list)
	router.GET("/strategy/grey", GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.get)
	router.POST("/strategy/grey", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyGrey), c.create)
	router.PUT("/strategy/grey", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.update)
	router.DELETE("/strategy/grey", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyGrey), c.del)
	router.PATCH("/strategy/grey/restore", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.restore)
	router.PATCH("/strategy/grey/stop", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.updateStop)
	router.GET("/strategy/grey/to-publishs", GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.toPublish)
	router.POST("/strategy/grey/publish", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyGrey), c.publish)
	router.POST("/strategy/grey/priority", GenAccessHandler(access.StrategyGreyEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.changePriority)
	router.GET("/strategy/grey/publish-history", GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.publishHistory)
}
