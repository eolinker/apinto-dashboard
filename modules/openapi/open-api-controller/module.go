package open_api_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type PluginDriver struct {
}

func (c *PluginDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *PluginDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewModule(info.Id, info.Name), nil
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}

type Module struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *Module) Middleware() []pm3.Middleware {
	return nil
}

func (c *Module) Name() string {
	return c.name
}

func (c *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewModule(id, name string) *Module {

	return &Module{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *Module) Apis() []pm3.Api {
	if !c.isInit {
		routers := AllRouters()
		apis := make([]pm3.Api, 0)
		for i := range routers {
			apis = append(apis, routers[i].Apis()...)
		}
		c.routers = apis
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}
