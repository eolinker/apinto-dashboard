package controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type emailDriver struct {
}

func (c *emailDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *emailDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewEmailModule(info.Id, info.Name), nil
}

func NewEmailDriver() apinto_module.Driver {
	return &emailDriver{}
}

type emailModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *emailModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *emailModule) Apis() []pm3.Api {
	return c.RoutersInfo()
}

func (c *emailModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *emailModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (c *emailModule) Name() string {
	return c.name
}

func NewEmailModule(id, name string) apinto_module.Module {

	return &emailModule{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *emailModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}

func (c *emailModule) initRouter() {
	emailCtl := newEmailController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/email",

			HandlerFunc: emailCtl.getEmail,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/email",

			HandlerFunc: emailCtl.createEmail,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/email",

			HandlerFunc: emailCtl.updateEmail,
		},
	}

}
