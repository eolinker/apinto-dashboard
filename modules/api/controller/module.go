package controller

import (
	"github.com/eolinker/apinto-dashboard/module"
)

type PluginDriver struct {
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}
func (c *PluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

// plugin
func (c *PluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewModule(name), nil
}

func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
	return "router/api"
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

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

type Module struct {
	isInit               bool
	name                 string
	routers              apinto_module.RoutersInfo
	filterOptionHandlers []apinto_module.IFilterOptionHandler
}

func (c *Module) Name() string {
	return c.name
}

func (c *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *Module) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *Module) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewModule(name string) *Module {
	m := &Module{name: name}
	m.initFilter()
	return m
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.routers = initRouter(c.name)
		c.isInit = true
	}
	return c.routers
}

func (c *Module) FilterOptionHandler() []apinto_module.IFilterOptionHandler {
	return c.filterOptionHandlers
}

func (c *Module) initFilter() {
	c.filterOptionHandlers = []apinto_module.IFilterOptionHandler{
		newFilterOption(),
	}
}
