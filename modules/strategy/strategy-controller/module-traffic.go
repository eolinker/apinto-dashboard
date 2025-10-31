package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyTrafficDriver struct {
}

func (c *StrategyTrafficDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *StrategyTrafficDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyTrafficModule(info.Id, info.Name), nil
}

func NewStrategyTraffic() apinto_module.Driver {
	return &StrategyTrafficDriver{}
}

type StrategyTrafficModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyTrafficModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyTrafficModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyTrafficModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyTrafficModule) Name() string {
	return c.name
}

func (c *StrategyTrafficModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyTrafficModule(id, name string) *StrategyTrafficModule {

	return &StrategyTrafficModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyTrafficModule) initRouter() {
	strategyTrafficController := newStrategyTrafficController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/traffic",

			HandlerFunc: strategyTrafficController.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/traffic",

			HandlerFunc: strategyTrafficController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/traffic",

			HandlerFunc: strategyTrafficController.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/traffic",

			HandlerFunc: strategyTrafficController.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/traffic",

			HandlerFunc: strategyTrafficController.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/traffic/restore",

			HandlerFunc: strategyTrafficController.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/traffic/enable",

			HandlerFunc: strategyTrafficController.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/traffic/disable",

			HandlerFunc: strategyTrafficController.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/traffic/to-publishs",

			HandlerFunc: strategyTrafficController.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/traffic/publish",

			HandlerFunc: strategyTrafficController.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/traffic/priority",

			HandlerFunc: strategyTrafficController.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/traffic/publish-history",

			HandlerFunc: strategyTrafficController.publishHistory,
		},
	}
}
