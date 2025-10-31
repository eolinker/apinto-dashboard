package application_controller

import (
	module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type PluginDriver struct {
}

func (p *PluginDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (p *PluginDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewModule(info.Id, info.Name), nil
}

func NewPluginDriver() module.Driver {
	return &PluginDriver{}
}

//
//
//func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
//	return "application"
//}

type Module struct {
	*pm3.ModuleTool

	isInit               bool
	name                 string
	routers              module.RoutersInfo
	filterOptionHandlers []module.IFilterOptionHandler
}

func (c *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *Module) Apis() []pm3.Api {
	if !c.isInit {
		c.routers = initRouter(c.name)
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}

func (c *Module) Middleware() []pm3.Middleware {
	return nil
}

func (c *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (c *Module) Name() string {
	return c.name
}

func NewModule(id, name string) *Module {
	m := &Module{
		ModuleTool: pm3.NewModuleTool(id, name),
		name:       name}
	m.initFilter()
	return m
}

func (c *Module) FilterOptionHandler() []module.IFilterOptionHandler {
	return c.filterOptionHandlers
}

func (c *Module) initFilter() {
	c.filterOptionHandlers = []module.IFilterOptionHandler{
		newFilterOption(),
	}
}
