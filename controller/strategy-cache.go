package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/gin-gonic/gin"
)

func RegisterStrategyCacheRouter(router gin.IRoutes) {
	strategyService := service.NewStrategyService(strategy_handler.NewStrategyCacheHandler("strategy-cache"), enum.StrategyCacheRuntimeKind)

	c := newStrategyController(strategyService)
	router.GET("/strategies/cache", genAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.list)
	router.GET("/strategy/cache", genAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.get)
	router.POST("/strategy/cache", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyCache), c.create)
	router.PUT("/strategy/cache", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.update)
	router.DELETE("/strategy/cache", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyCache), c.del)
	router.PATCH("/strategy/cache/restore", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.restore)
	router.PATCH("/strategy/cache/stop", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.updateStop)
	router.GET("/strategy/cache/to-publishs", genAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.toPublish)
	router.POST("/strategy/cache/publish", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindStrategyCache), c.publish)
	router.POST("/strategy/cache/priority", genAccessHandler(access.StrategyCacheEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.changePriority)
	router.GET("/strategy/cache/publish-history", genAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.publishHistory)
}
