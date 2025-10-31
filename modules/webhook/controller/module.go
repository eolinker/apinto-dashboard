package controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type webhookDriver struct {
}

func (c *webhookDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *webhookDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewWebhookModule(info.Id, info.Name), nil
}

func NewWebhookDriver() apinto_module.Driver {
	return &webhookDriver{}
}

//func (c *webhookDriver) GetPluginFrontend(moduleName string) string {
//	return "system/webhook"
//}

type webhookModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *webhookModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *webhookModule) Apis() []pm3.Api {

	return c.routers
}

func (c *webhookModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *webhookModule) Name() string {
	return c.name
}

func (c *webhookModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewWebhookModule(id, name string) apinto_module.Module {

	m := &webhookModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
	m.initRouter()
	return m
}

func (c *webhookModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *webhookModule) initRouter() {
	webhookCtl := newWebhookController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/webhooks",

			HandlerFunc: webhookCtl.webhooks,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/webhook",

			HandlerFunc: webhookCtl.webhook,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/webhook",

			HandlerFunc: webhookCtl.createWebhook,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/webhook",

			HandlerFunc: webhookCtl.updateWebhook,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/webhook",

			HandlerFunc: webhookCtl.delWebhook,
		},
	}
	c.InitAccess(c.routers)

}
