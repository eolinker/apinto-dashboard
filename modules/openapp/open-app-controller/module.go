package open_app_controller

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

// plugin
func (c *PluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewModule(name), nil
}

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
	return "system/ext-app"
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
	return true
}

type Module struct {
	isInit     bool
	name       string
	routers    apinto_module.RoutersInfo
	middleware []apinto_module.MiddlewareHandler
}

func (c *Module) MiddlewaresInfo() []apinto_module.MiddlewareHandler {
	return c.middleware
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
	return c, true
}

func NewModule(name string) *Module {

	return &Module{name: name}
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.routers, c.middleware = initRouter(c.name)

		c.isInit = true
	}
	return c.routers
}
