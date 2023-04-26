package discovery_controller

import (
	"github.com/eolinker/apinto-module"
)

type PluginDriver struct {
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}
func (c *PluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
	return "upstream/discovery"
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
