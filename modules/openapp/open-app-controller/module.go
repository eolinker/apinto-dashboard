package open_app_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type PluginDriver struct {
}

func (c *PluginDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {

	ms, acs, fs, err = pm3.ReadPluginAssembly(info)

	return
}

func (c *PluginDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewModule(info.Id, info.Name), nil
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}

type Module struct {
	*pm3.ModuleTool
	routers    apinto_module.RoutersInfo
	middleware []apinto_module.MiddlewareHandler
}

func (c *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *Module) Apis() []pm3.Api {
	return c.RoutersInfo()
}

func (c *Module) Middleware() []pm3.Middleware {
	return c.middleware
}

func (c *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewModule(id, name string) *Module {
	ModuleTool := pm3.NewModuleTool(id, name)

	routers, middleware := initRouter(name)
	ModuleTool.InitAccess(routers)
	return &Module{
		ModuleTool: ModuleTool,
		routers:    routers,
		middleware: middleware,
	}
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {

	return c.routers
}
