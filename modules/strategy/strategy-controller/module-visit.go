package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyVisitDriver struct {
}

func (c *StrategyVisitDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *StrategyVisitDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyVisitModule(info.Id, info.Name), nil

}

func NewStrategyVisit() apinto_module.Driver {
	return &StrategyVisitDriver{}
}

type StrategyVisitModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyVisitModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyVisitModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyVisitModule) Name() string {
	return c.name
}

func (c *StrategyVisitModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyVisitModule(id, name string) *StrategyVisitModule {

	return &StrategyVisitModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyVisitModule) Apis() []pm3.Api {

	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyVisitModule) initRouter() {
	strategyVisitController := newStrategyVisitController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/visit",

			HandlerFunc: strategyVisitController.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/visit",

			HandlerFunc: strategyVisitController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/visit",

			HandlerFunc: strategyVisitController.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/visit",

			HandlerFunc: strategyVisitController.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/visit",

			HandlerFunc: strategyVisitController.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/visit/restore",

			HandlerFunc: strategyVisitController.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/visit/enable",

			HandlerFunc: strategyVisitController.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/visit/disable",

			HandlerFunc: strategyVisitController.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/visit/to-publishs",

			HandlerFunc: strategyVisitController.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/visit/publish",

			HandlerFunc: strategyVisitController.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/visit/priority",

			HandlerFunc: strategyVisitController.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/visit/publish-history",

			HandlerFunc: strategyVisitController.publishHistory,
		},
	}
}
