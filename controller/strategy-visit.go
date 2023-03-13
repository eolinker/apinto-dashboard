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
	router.GET("/strategies/visit", genAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.list)
	router.GET("/strategy/visit", genAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.get)
	router.POST("/strategy/visit", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyVisit), c.create)
	router.PUT("/strategy/visit", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.update)
	router.DELETE("/strategy/visit", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyVisit), c.del)
	router.PATCH("/strategy/visit/restore", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.restore)
	router.PATCH("/strategy/visit/stop", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.updateStop)
	router.GET("/strategy/visit/to-publishs", genAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.toPublish)
	router.POST("/strategy/visit/publish", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindStrategyVisit), c.publish)
	router.POST("/strategy/visit/priority", genAccessHandler(access.StrategyVisitEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyVisit), c.changePriority)
	router.GET("/strategy/visit/publish-history", genAccessHandler(access.StrategyVisitView, access.StrategyVisitEdit), c.publishHistory)
}
