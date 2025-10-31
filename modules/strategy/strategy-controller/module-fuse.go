package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyFuseDriver struct {
}

func (c *StrategyFuseDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *StrategyFuseDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyFuseModule(info.Id, info.Name), nil
}

func NewStrategyFuse() apinto_module.Driver {
	return &StrategyFuseDriver{}
}

//
//func (c *StrategyFuseDriver) GetPluginFrontend(moduleName string) string {
//	return "serv-governance/fuse"
//}
//

type StrategyFuseModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyFuseModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyFuseModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyFuseModule) Name() string {
	return c.name
}

func (c *StrategyFuseModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyFuseModule(id, name string) *StrategyFuseModule {

	return &StrategyFuseModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyFuseModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyFuseModule) initRouter() {
	strategyFuseController := newStrategyFuseController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/fuse",

			HandlerFunc: strategyFuseController.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/fuse",

			HandlerFunc: strategyFuseController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/fuse",

			HandlerFunc: strategyFuseController.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/fuse",

			HandlerFunc: strategyFuseController.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/fuse",

			HandlerFunc: strategyFuseController.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/fuse/restore",

			HandlerFunc: strategyFuseController.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/fuse/enable",

			HandlerFunc: strategyFuseController.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/fuse/disable",

			HandlerFunc: strategyFuseController.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/fuse/to-publishs",

			HandlerFunc: strategyFuseController.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/fuse/publish",

			HandlerFunc: strategyFuseController.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/fuse/priority",

			HandlerFunc: strategyFuseController.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/fuse/publish-history",

			HandlerFunc: strategyFuseController.publishHistory,
		},
	}
}
