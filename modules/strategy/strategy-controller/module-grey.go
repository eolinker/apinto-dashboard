package strategy_controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type StrategyGreyDriver struct {
}

func NewStrategyGrey() apinto_module.Driver {
	return &StrategyGreyDriver{}
}

func (c *StrategyGreyDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewStrategyGreyModule(name), nil
}

func (c *StrategyGreyDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *StrategyGreyDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type StrategyGreyModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyGreyModule) Name() string {
	return c.name
}

func (c *StrategyGreyModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *StrategyGreyModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *StrategyGreyModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewStrategyGreyModule(name string) *StrategyGreyModule {

	return &StrategyGreyModule{name: name}
}

func (c *StrategyGreyModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *StrategyGreyModule) initRouter() {
	strategyGreyController := newStrategyGreyController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategies/grey",
			Handler:     "strategy-grey.list",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/grey",
			Handler:     "strategy-grey.get",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/grey",
			Handler:     "strategy-grey.create",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/strategy/grey",
			Handler:     "strategy-grey.update",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/strategy/grey",
			Handler:     "strategy-grey.del",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.del},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/grey/restore",
			Handler:     "strategy-grey.restore",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.restore},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/grey/stop",
			Handler:     "strategy-grey.updateStop",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.updateStop},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/grey/to-publishs",
			Handler:     "strategy-grey.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.toPublish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/grey/publish",
			Handler:     "strategy-grey.publish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.publish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/grey/priority",
			Handler:     "strategy-grey.changePriority",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.changePriority},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/grey/publish-history",
			Handler:     "strategy-grey.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{strategyGreyController.publishHistory},
		},
	}
}
