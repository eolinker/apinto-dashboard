package controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type webhookDriver struct {
}

func NewWebhookDriver() apinto_module.Driver {
	return &webhookDriver{}
}

func (c *webhookDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewWebhookModule(name), nil
}

func (c *webhookDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *webhookDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *webhookDriver) GetPluginFrontend(moduleName string) string {
	return "system/webhook"
}

func (c *webhookDriver) IsPluginVisible() bool {
	return true
}

func (c *webhookDriver) IsShowServer() bool {
	return false
}

func (c *webhookDriver) IsCanUninstall() bool {
	return false
}

func (c *webhookDriver) IsCanDisable() bool {
	return true
}

type webhookModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *webhookModule) Name() string {
	return c.name
}

func (c *webhookModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *webhookModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *webhookModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewWebhookModule(name string) apinto_module.Module {

	return &webhookModule{name: name}
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
			Method:      http.MethodGet,
			Path:        "/api/webhooks",
			Handler:     "webhook.getList",
			HandlerFunc: []apinto_module.HandlerFunc{webhookCtl.webhooks},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/webhook",
			Handler:     "webhook.getInfo",
			HandlerFunc: []apinto_module.HandlerFunc{webhookCtl.webhook},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/webhook",
			Handler:     "webhook.createWebhook",
			HandlerFunc: []apinto_module.HandlerFunc{webhookCtl.createWebhook},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/webhook",
			Handler:     "webhook.alterWebhook",
			HandlerFunc: []apinto_module.HandlerFunc{webhookCtl.updateWebhook},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/webhook",
			Handler:     "webhook.alterWebhook",
			HandlerFunc: []apinto_module.HandlerFunc{webhookCtl.delWebhook},
		},
	}
}
