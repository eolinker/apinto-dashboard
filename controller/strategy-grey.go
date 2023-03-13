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
	router.GET("/strategies/grey", genAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.list)
	router.GET("/strategy/grey", genAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.get)
	router.POST("/strategy/grey", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyGrey), c.create)
	router.PUT("/strategy/grey", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.update)
	router.DELETE("/strategy/grey", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyGrey), c.del)
	router.PATCH("/strategy/grey/restore", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.restore)
	router.PATCH("/strategy/grey/stop", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.updateStop)
	router.GET("/strategy/grey/to-publishs", genAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.toPublish)
	router.POST("/strategy/grey/publish", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindStrategyGrey), c.publish)
	router.POST("/strategy/grey/priority", genAccessHandler(access.StrategyGreyEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.changePriority)
	router.GET("/strategy/grey/publish-history", genAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.publishHistory)
}
