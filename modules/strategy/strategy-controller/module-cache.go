package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyCacheDriver struct {
}

func (c *StrategyCacheDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *StrategyCacheDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyCacheModule(info.Id, info.Name), nil
}

func NewStrategyCache() apinto_module.Driver {
	return &StrategyCacheDriver{}
}

//func (c *StrategyCacheDriver) GetPluginFrontend(moduleName string) string {
//	return "serv-governance/cache"
//}

type StrategyCacheModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyCacheModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyCacheModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyCacheModule) Name() string {
	return c.name
}

func (c *StrategyCacheModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyCacheModule(id, name string) *StrategyCacheModule {

	return &StrategyCacheModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyCacheModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyCacheModule) initRouter() {
	strategyCacheController := newStrategyCacheController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/cache",

			HandlerFunc: strategyCacheController.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/cache",

			HandlerFunc: strategyCacheController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/cache",

			HandlerFunc: strategyCacheController.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/cache",

			HandlerFunc: strategyCacheController.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/cache",

			HandlerFunc: strategyCacheController.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/cache/restore",

			HandlerFunc: strategyCacheController.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/cache/enable",

			HandlerFunc: strategyCacheController.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/cache/disable",

			HandlerFunc: strategyCacheController.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/cache/to-publishs",

			HandlerFunc: strategyCacheController.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/cache/publish",

			HandlerFunc: strategyCacheController.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/cache/priority",

			HandlerFunc: strategyCacheController.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/cache/publish-history",

			HandlerFunc: strategyCacheController.publishHistory,
		},
	}
}
