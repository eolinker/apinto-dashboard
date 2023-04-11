package plugin_template_controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type PluginTemplateDriver struct {
}

func NewPluginTemplateDriver() apinto_module.Driver {
	return &PluginTemplateDriver{}
}

func (c *PluginTemplateDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewPluginTemplateModule(name), nil
}

func (c *PluginTemplateDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *PluginTemplateDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type PluginTemplateModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *PluginTemplateModule) Name() string {
	return c.name
}

func (c *PluginTemplateModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *PluginTemplateModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *PluginTemplateModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewPluginTemplateModule(name string) *PluginTemplateModule {

	return &PluginTemplateModule{name: name}
}

func (c *PluginTemplateModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *PluginTemplateModule) initRouter() {
	pluginTemplateCtl := newPluginTemplateController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/templates",
			Handler:     "plugin-template.templates",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.templates},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/template/enum",
			Handler:     "plugin-template.templateEnum",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.templateEnum},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/plugin/template",
			Handler:     "plugin-template.createTemplate",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.createTemplate},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/plugin/template",
			Handler:     "plugin-template.updateTemplate",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.updateTemplate},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/plugin/template",
			Handler:     "plugin-template.delTemplate",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.delTemplate},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/template",
			Handler:     "plugin-template.template",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.template},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/template/onlines",
			Handler:     "plugin-template.onlines",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.onlines},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/plugin/template/online",
			Handler:     "plugin-template.online",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.online},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/plugin/template/offline",
			Handler:     "plugin-template.offline",
			HandlerFunc: []apinto_module.HandlerFunc{pluginTemplateCtl.offline},
		},
	}
}
