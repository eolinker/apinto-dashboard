package plugin_template_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type PluginTemplateDriver struct {
}

func (c *PluginTemplateDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *PluginTemplateDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewPluginTemplateModule(info.Id, info.Name), nil
}

func NewPluginTemplateDriver() apinto_module.Driver {
	return &PluginTemplateDriver{}
}

type PluginTemplateModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *PluginTemplateModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *PluginTemplateModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}

func (c *PluginTemplateModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *PluginTemplateModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (c *PluginTemplateModule) Name() string {
	return c.name
}

func NewPluginTemplateModule(id, name string) *PluginTemplateModule {

	return &PluginTemplateModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *PluginTemplateModule) initRouter() {
	pluginTemplateCtl := newPluginTemplateController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/templates",

			HandlerFunc: pluginTemplateCtl.templates,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/template/enum",

			HandlerFunc: pluginTemplateCtl.templateEnum,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/plugin/template",

			HandlerFunc: pluginTemplateCtl.createTemplate,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/plugin/template",

			HandlerFunc: pluginTemplateCtl.updateTemplate,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/plugin/template",

			HandlerFunc: pluginTemplateCtl.delTemplate,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/template",

			HandlerFunc: pluginTemplateCtl.template,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/plugin/template/onlines",

			HandlerFunc: pluginTemplateCtl.onlines,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/plugin/template/online",

			HandlerFunc: pluginTemplateCtl.online,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/plugin/template/offline",

			HandlerFunc: pluginTemplateCtl.offline,
		},
	}
}
