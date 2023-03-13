package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyTrafficRouter(router gin.IRoutes) {
	strategyService := service.NewStrategyService(strategy_handler.NewStrategyTrafficHandler("strategy-limiting"), enum.StrategyTrafficRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/traffic", GenAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.list)
	router.GET("/strategy/traffic", GenAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.get)
	router.POST("/strategy/traffic", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyTraffic), c.create)
	router.PUT("/strategy/traffic", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.update)
	router.DELETE("/strategy/traffic", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyTraffic), c.del)
	router.PATCH("/strategy/traffic/restore", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.restore)
	router.PATCH("/strategy/traffic/stop", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.updateStop)
	router.GET("/strategy/traffic/to-publishs", GenAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.toPublish)
	router.POST("/strategy/traffic/publish", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyTraffic), c.publish)
	router.POST("/strategy/traffic/priority", GenAccessHandler(access.StrategyTrafficEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.changePriority)
	router.GET("/strategy/traffic/publish-history", GenAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.publishHistory)
}
