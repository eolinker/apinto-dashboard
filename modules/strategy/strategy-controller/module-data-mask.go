package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type StrategyDataMaskDriver struct {
}

func (c *StrategyDataMaskDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *StrategyDataMaskDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewStrategyDataMaskDriverModule(info.Id, info.Name), nil

}

func NewStrategyDataMask() apinto_module.Driver {
	return &StrategyDataMaskDriver{}
}

type StrategyDataMaskModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyDataMaskModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *StrategyDataMaskModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *StrategyDataMaskModule) Name() string {
	return c.name
}

func (c *StrategyDataMaskModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewStrategyDataMaskDriverModule(id, name string) *StrategyDataMaskModule {

	return &StrategyDataMaskModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *StrategyDataMaskModule) Apis() []pm3.Api {

	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *StrategyDataMaskModule) initRouter() {
	controller := newStrategyDataMaskController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/strategies/data-mask",

			HandlerFunc: controller.list,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/data-mask",

			HandlerFunc: controller.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/data-mask",

			HandlerFunc: controller.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/strategy/data-mask",

			HandlerFunc: controller.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/strategy/data-mask",

			HandlerFunc: controller.del,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/data-mask/restore",

			HandlerFunc: controller.restore,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/data-mask/enable",

			HandlerFunc: controller.enable,
		},
		{
			Method: http.MethodPatch,
			Path:   "/api/strategy/data-mask/disable",

			HandlerFunc: controller.disable,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/data-mask/to-publishs",

			HandlerFunc: controller.toPublish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/data-mask/publish",

			HandlerFunc: controller.publish,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/strategy/data-mask/priority",

			HandlerFunc: controller.changePriority,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/strategy/data-mask/publish-history",

			HandlerFunc: controller.publishHistory,
		},
	}
}
