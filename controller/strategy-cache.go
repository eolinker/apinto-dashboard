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
	router.GET("/strategies/cache", GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.list)
	router.GET("/strategy/cache", GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.get)
	router.POST("/strategy/cache", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindStrategyCache), c.create)
	router.PUT("/strategy/cache", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.update)
	router.DELETE("/strategy/cache", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindStrategyCache), c.del)
	router.PATCH("/strategy/cache/restore", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.restore)
	router.PATCH("/strategy/cache/stop", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.updateStop)
	router.GET("/strategy/cache/to-publishs", GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.toPublish)
	router.POST("/strategy/cache/publish", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindStrategyCache), c.publish)
	router.POST("/strategy/cache/priority", GenAccessHandler(access.StrategyCacheEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindStrategyCache), c.changePriority)
	router.GET("/strategy/cache/publish-history", GenAccessHandler(access.StrategyCacheView, access.StrategyCacheEdit), c.publishHistory)
}
