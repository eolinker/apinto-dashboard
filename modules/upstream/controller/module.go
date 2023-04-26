package upstream_controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-module"
	"net/http"
)

type UpstreamDriver struct {
}

func NewUpstreamDriver() apinto_module.Driver {
	return &UpstreamDriver{}
}

func (c *UpstreamDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewUpstreamModule(name), nil
}

func (c *UpstreamDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *UpstreamDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *UpstreamDriver) GetPluginFrontend(moduleName string) string {
	return "upstream/upstream"
}

func (c *UpstreamDriver) IsPluginVisible() bool {
	return true
}

func (c *UpstreamDriver) IsShowServer() bool {
	return false
}

func (c *UpstreamDriver) IsCanUninstall() bool {
	return false
}

func (c *UpstreamDriver) IsCanDisable() bool {
	return false
}

type UpstreamModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *UpstreamModule) Name() string {
	return c.name
}

func (c *UpstreamModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *UpstreamModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *UpstreamModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewUpstreamModule(name string) *UpstreamModule {

	return &UpstreamModule{name: name}
}

func (c *UpstreamModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *UpstreamModule) initRouter() {
	upstreamCtl := newUpstreamController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/services",
			Handler:     "upstream.getList",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.getList},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/service",
			Handler:     "upstream.getInfo",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.getInfo},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/service",
			Handler:     "upstream.alter",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.alter},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/service",
			Handler:     "upstream.create",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.create},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/service",
			Handler:     "upstream.del",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.del},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/service/enum",
			Handler:     "upstream.getEnum",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.getEnum},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/service/:service_name/online",
			Handler:     "upstream.online",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.online},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/service/:service_name/offline",
			Handler:     "upstream.offline",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, upstreamCtl.offline},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/service/:service_name/onlines",
			Handler:     "upstream.getOnlineList",
			HandlerFunc: []apinto_module.HandlerFunc{upstreamCtl.getOnlineList},
		},
	}
}
