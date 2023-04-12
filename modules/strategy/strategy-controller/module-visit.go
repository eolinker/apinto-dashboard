package strategy_controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type StrategyVisitDriver struct {
}

func NewStrategyVisit() apinto_module.Driver {
	return &StrategyVisitDriver{}
}

func (c *StrategyVisitDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewStrategyVisitModule(name), nil
}

func (c *StrategyVisitDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *StrategyVisitDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type StrategyVisitModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyVisitModule) Name() string {
	return c.name
}

func (c *StrategyVisitModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *StrategyVisitModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *StrategyVisitModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewStrategyVisitModule(name string) *StrategyVisitModule {

	return &StrategyVisitModule{name: name}
}

func (c *StrategyVisitModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *StrategyVisitModule) initRouter() {
	strategyVisitController := newStrategyVisitController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategies/visit",
			Handler:     "strategy-visit.list",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/visit",
			Handler:     "strategy-visit.get",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/visit",
			Handler:     "strategy-visit.create",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/strategy/visit",
			Handler:     "strategy-visit.update",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/strategy/visit",
			Handler:     "strategy-visit.del",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.del},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/visit/restore",
			Handler:     "strategy-visit.restore",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.restore},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/visit/stop",
			Handler:     "strategy-visit.updateStop",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.updateStop},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/visit/to-publishs",
			Handler:     "strategy-visit.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.toPublish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/visit/publish",
			Handler:     "strategy-visit.publish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.publish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/visit/priority",
			Handler:     "strategy-visit.changePriority",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.changePriority},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/visit/publish-history",
			Handler:     "strategy-visit.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{strategyVisitController.publishHistory},
		},
	}
}
