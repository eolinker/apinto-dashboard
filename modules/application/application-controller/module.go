package application_controller

import (
	module "github.com/eolinker/apinto-dashboard/module"
)

type PluginDriver struct {
}

func NewPluginDriver() module.Driver {
	return &PluginDriver{}
}
func (c *PluginDriver) CreatePlugin(define interface{}) (module.Plugin, error) {
	return c, nil
}

func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
	return "application"
}

func (c *PluginDriver) IsPluginVisible() bool {
	return true
}

func (c *PluginDriver) IsShowServer() bool {
	return false
}

func (c *PluginDriver) IsCanUninstall() bool {
	return false
}

func (c *PluginDriver) IsCanDisable() bool {
	return false
}

// plugin
func (c *PluginDriver) CreateModule(name string, config interface{}) (module.Module, error) {
	return NewModule(name), nil
}

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

type Module struct {
	isInit               bool
	name                 string
	routers              module.RoutersInfo
	filterOptionHandlers []module.IFilterOptionHandler
}

func (c *Module) Name() string {
	return c.name
}

func (c *Module) Support() (module.ProviderSupport, bool) {
	return nil, false
}

func (c *Module) Routers() (module.Routers, bool) {
	return c, true
}

func (c *Module) Middleware() (module.Middleware, bool) {
	return nil, false
}

func NewModule(name string) *Module {
	m := &Module{name: name}
	m.initFilter()
	return m
}

func (c *Module) RoutersInfo() module.RoutersInfo {
	if !c.isInit {
		c.routers = initRouter(c.name)
		c.isInit = true
	}
	return c.routers
}

func (c *Module) FilterOptionHandler() []module.IFilterOptionHandler {
	return c.filterOptionHandlers
}

func (c *Module) initFilter() {
	c.filterOptionHandlers = []module.IFilterOptionHandler{
		newFilterOption(),
	}
}
