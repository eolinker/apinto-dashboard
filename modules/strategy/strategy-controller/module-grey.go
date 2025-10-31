package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyGreyDriver struct {
}

func (s *StrategyGreyDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (s *StrategyGreyDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyGreyModule(info.Id, info.Name), nil
}

func NewStrategyGrey() apinto_module.Driver {
	return &StrategyGreyDriver{}
}

//
//func (c *StrategyGreyDriver) GetPluginFrontend(moduleName string) string {
//	return "serv-governance/grey"
//}

type StrategyGreyModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyGreyModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyGreyModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyGreyModule) Name() string {
	return c.name
}

func (c *StrategyGreyModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyGreyModule(id, name string) *StrategyGreyModule {

	return &StrategyGreyModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyGreyModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyGreyModule) initRouter() {
	strategyGreyController := newStrategyGreyController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/grey",

			HandlerFunc: strategyGreyController.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/grey",

			HandlerFunc: strategyGreyController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/grey",

			HandlerFunc: strategyGreyController.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/grey",

			HandlerFunc: strategyGreyController.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/grey",

			HandlerFunc: strategyGreyController.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/grey/restore",

			HandlerFunc: strategyGreyController.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/grey/enable",

			HandlerFunc: strategyGreyController.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/grey/disable",

			HandlerFunc: strategyGreyController.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/grey/to-publishs",

			HandlerFunc: strategyGreyController.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/grey/publish",

			HandlerFunc: strategyGreyController.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/grey/priority",

			HandlerFunc: strategyGreyController.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/grey/publish-history",

			HandlerFunc: strategyGreyController.publishHistory,
		},
	}
}
