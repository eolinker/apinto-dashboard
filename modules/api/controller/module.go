package controller

import (
	"github.com/eolinker/apinto-module"
)

type PluginDriver struct {
}

func (c *PluginDriver) GetPluginFrontend(define interface{}) string {
	//TODO implement me
	panic("implement me")
}

func (c *PluginDriver) IsPluginInvisible(define interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (c *PluginDriver) IsShowServer(define interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (c *PluginDriver) IsCanUninstall(define interface{}) bool {
	//TODO implement me
	panic("implement me")
}

func (c *PluginDriver) IsCanDisable(define interface{}) bool {
	//TODO implement me
	panic("implement me")
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

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

type Module struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
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

	return &Module{name: name}
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.routers = initRouter(c.name)
		c.isInit = true
	}
	return c.routers
}
