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
	router.GET("/strategies/traffic", genAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.list)
	router.GET("/strategy/traffic", genAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.get)
	router.POST("/strategy/traffic", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyTraffic), c.create)
	router.PUT("/strategy/traffic", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.update)
	router.DELETE("/strategy/traffic", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyTraffic), c.del)
	router.PATCH("/strategy/traffic/restore", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.restore)
	router.PATCH("/strategy/traffic/stop", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.updateStop)
	router.GET("/strategy/traffic/to-publishs", genAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.toPublish)
	router.POST("/strategy/traffic/publish", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindStrategyTraffic), c.publish)
	router.POST("/strategy/traffic/priority", genAccessHandler(access.StrategyTrafficEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyTraffic), c.changePriority)
	router.GET("/strategy/traffic/publish-history", genAccessHandler(access.StrategyTrafficView, access.StrategyTrafficEdit), c.publishHistory)
}
