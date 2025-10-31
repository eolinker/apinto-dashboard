package plugin_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type PluginDriver struct {
}

func (c *PluginDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *PluginDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewPluginModule(info.Id, info.Name), nil
}

func NewPluginDriver() apinto_module.Driver {
	return &PluginDriver{}
}

//
//
//func (c *PluginDriver) GetPluginFrontend(moduleName string) string {
//	return "deploy/plugin"
//}

type PluginModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *PluginModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *PluginModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *PluginModule) Name() string {
	return c.name
}

func (c *PluginModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewPluginModule(id, name string) *PluginModule {

	return &PluginModule{name: name, ModuleTool: pm3.NewModuleTool(id, name)}
}

func (c *PluginModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}

func (c *PluginModule) initRouter() {
	pluginCtl := newPluginController()
	clusterPluginCtl := newPluginClusterController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/plugins",

			HandlerFunc: pluginCtl.plugins,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/basic/info/plugins",

			HandlerFunc: pluginCtl.basicInfoPlugins,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin",

			HandlerFunc: pluginCtl.plugin,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/plugin",

			HandlerFunc: pluginCtl.createPlugin,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/plugin",

			HandlerFunc: pluginCtl.updatePlugin,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/plugin",

			HandlerFunc: pluginCtl.delPlugin,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/extendeds",

			HandlerFunc: pluginCtl.pluginExtendeds,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugins/render",

			HandlerFunc: pluginCtl.pluginRender,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/plugin/sort",

			HandlerFunc: pluginCtl.pluginSort,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/enum",

			HandlerFunc: pluginCtl.pluginEnum,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/plugins",

			HandlerFunc: clusterPluginCtl.plugins,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/plugin",

			HandlerFunc: clusterPluginCtl.getPlugin,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/plugin",

			HandlerFunc: clusterPluginCtl.editPlugin,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/plugin/publish",

			HandlerFunc: clusterPluginCtl.publish,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/plugin/to-publish",

			HandlerFunc: clusterPluginCtl.toPublish,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/plugin/publish-history",

			HandlerFunc: clusterPluginCtl.publishHistory,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/plugin/update-history",

			HandlerFunc: clusterPluginCtl.updateHistory,
		},
	}
}
