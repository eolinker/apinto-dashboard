package plugin_controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type PluginDriver struct {
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}

func (c *PluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewPluginModule(name), nil
}

func (c *PluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *PluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type PluginModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *PluginModule) Name() string {
	return c.name
}

func (c *PluginModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *PluginModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *PluginModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewPluginModule(name string) *PluginModule {

	return &PluginModule{name: name}
}

func (c *PluginModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *PluginModule) initRouter() {
	pluginCtl := newPluginController()
	clusterPluginCtl := newPluginClusterController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/plugins",
			Handler:     "plugin.plugins",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.plugins},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/basic/info/plugins",
			Handler:     "plugin.basicInfoPlugins",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.basicInfoPlugins},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin",
			Handler:     "plugin.plugin",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.plugin},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/plugin",
			Handler:     "plugin.createPlugin",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.createPlugin},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/plugin",
			Handler:     "plugin.updatePlugin",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.updatePlugin},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/plugin",
			Handler:     "plugin.delPlugin",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.delPlugin},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/extendeds",
			Handler:     "plugin.pluginExtendeds",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.pluginExtendeds},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugins/render",
			Handler:     "plugin.pluginRender",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.pluginRender},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/plugin/sort",
			Handler:     "plugin.pluginSort",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.pluginSort},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/enum",
			Handler:     "plugin.pluginEnum",
			HandlerFunc: []apinto_module.HandlerFunc{pluginCtl.pluginEnum},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/plugins",
			Handler:     "cluster-plugin.plugins",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.plugins},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/plugin",
			Handler:     "cluster-plugin.getPlugin",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.getPlugin},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/plugin",
			Handler:     "cluster-plugin.editPlugin",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.editPlugin},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/plugin/publish",
			Handler:     "cluster-plugin.publish",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.publish},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/plugin/to-publish",
			Handler:     "cluster-plugin.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.toPublish},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/plugin/publish-history",
			Handler:     "cluster-plugin.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.publishHistory},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/plugin/update-history",
			Handler:     "cluster-plugin.updateHistory",
			HandlerFunc: []apinto_module.HandlerFunc{clusterPluginCtl.updateHistory},
		},
	}
}
