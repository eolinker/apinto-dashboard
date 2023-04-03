package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyGreyRouter(router gin.IRoutes) {
	strategyService := strategy_service.NewStrategyService(strategy_handler.NewStrategyGreyHandler("strategy-grey"), enum.StrategyGreyRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/grey", controller.GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.list)
	router.GET("/strategy/grey", controller.GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.get)
	router.POST("/strategy/grey", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyGrey), c.create)
	router.PUT("/strategy/grey", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.update)
	router.DELETE("/strategy/grey", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyGrey), c.del)
	router.PATCH("/strategy/grey/restore", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.restore)
	router.PATCH("/strategy/grey/stop", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.updateStop)
	router.GET("/strategy/grey/to-publishs", controller.GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.toPublish)
	router.POST("/strategy/grey/publish", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyGrey), c.publish)
	router.POST("/strategy/grey/priority", controller.GenAccessHandler(access.StrategyGreyEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyGrey), c.changePriority)
	router.GET("/strategy/grey/publish-history", controller.GenAccessHandler(access.StrategyGreyView, access.StrategyGreyEdit), c.publishHistory)
}
