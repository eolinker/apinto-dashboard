package strategy_controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-module"
	"net/http"
)

type StrategyCacheDriver struct {
}

func NewStrategyCache() apinto_module.Driver {
	return &StrategyCacheDriver{}
}

func (c *StrategyCacheDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewStrategyCacheModule(name), nil
}

func (c *StrategyCacheDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *StrategyCacheDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type StrategyCacheModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyCacheModule) Name() string {
	return c.name
}

func (c *StrategyCacheModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *StrategyCacheModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *StrategyCacheModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewStrategyCacheModule(name string) *StrategyCacheModule {

	return &StrategyCacheModule{name: name}
}

func (c *StrategyCacheModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *StrategyCacheModule) initRouter() {
	strategyCacheController := newStrategyCacheController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategies/cache",
			Handler:     "strategy-cache.list",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/cache",
			Handler:     "strategy-cache.get",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/cache",
			Handler:     "strategy-cache.create",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/strategy/cache",
			Handler:     "strategy-cache.update",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/strategy/cache",
			Handler:     "strategy-cache.del",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.del},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/cache/restore",
			Handler:     "strategy-cache.restore",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypeEdit.Handler, strategyCacheController.restore},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/cache/stop",
			Handler:     "strategy-cache.updateStop",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypeEdit.Handler, strategyCacheController.updateStop},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/cache/to-publishs",
			Handler:     "strategy-cache.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.toPublish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/cache/publish",
			Handler:     "strategy-cache.publish",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, strategyCacheController.publish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/cache/priority",
			Handler:     "strategy-cache.changePriority",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypeEdit.Handler, strategyCacheController.changePriority},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/cache/publish-history",
			Handler:     "strategy-cache.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{strategyCacheController.publishHistory},
		},
	}
}
