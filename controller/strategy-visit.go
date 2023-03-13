package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyVisitRouter(router gin.IRoutes) {
	strategyService := service.NewStrategyService(strategy_handler.NewStrategyVisitHandler("strategy-visit"), enum.StrategyVisitRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/visit", GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.list)
	router.GET("/strategy/visit", GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.get)
	router.POST("/strategy/visit", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyVisit), c.create)
	router.PUT("/strategy/visit", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.update)
	router.DELETE("/strategy/visit", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyVisit), c.del)
	router.PATCH("/strategy/visit/restore", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.restore)
	router.PATCH("/strategy/visit/stop", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.updateStop)
	router.GET("/strategy/visit/to-publishs", GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.toPublish)
	router.POST("/strategy/visit/publish", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyVisit), c.publish)
	router.POST("/strategy/visit/priority", GenAccessHandler(access.StrategyVisitEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.changePriority)
	router.GET("/strategy/visit/publish-history", GenAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.publishHistory)
}
