package controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"net/http"
)

type emailDriver struct {
}

func NewEmailDriver() apinto_module.Driver {
	return &emailDriver{}
}

func (c *emailDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewEmailModule(name), nil
}

func (c *emailDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *emailDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *emailDriver) GetPluginFrontend(moduleName string) string {
	return "system/email"
}

func (c *emailDriver) IsPluginVisible() bool {
	return true
}

func (c *emailDriver) IsShowServer() bool {
	return false
}

func (c *emailDriver) IsCanUninstall() bool {
	return false
}

func (c *emailDriver) IsCanDisable() bool {
	return true
}

type emailModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *emailModule) Name() string {
	return c.name
}

func (c *emailModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *emailModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *emailModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewEmailModule(name string) apinto_module.Module {

	return &emailModule{name: name}
}

func (c *emailModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *emailModule) initRouter() {
	emailCtl := newEmailController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/email",
			Handler:     "email.getInfo",
			HandlerFunc: []apinto_module.HandlerFunc{emailCtl.getEmail},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/email",
			Handler:     "email.createEmail",
			HandlerFunc: []apinto_module.HandlerFunc{emailCtl.createEmail},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/email",
			Handler:     "email.alterEmail",
			HandlerFunc: []apinto_module.HandlerFunc{emailCtl.updateEmail},
		},
	}
}
