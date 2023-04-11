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
	return NewClusterModule(name), nil
}

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

type ClusterModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *ClusterModule) Name() string {
	return c.name
}

func (c *ClusterModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *ClusterModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *ClusterModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewClusterModule(name string) *ClusterModule {

	return &ClusterModule{name: name}
}

func (c *ClusterModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.routers = initRouter(c.name)
		c.isInit = true
	}
	return c.routers
}
