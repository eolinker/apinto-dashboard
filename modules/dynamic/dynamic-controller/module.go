package dynamic_controller

import (
	"fmt"
	"net/http"

	"github.com/eolinker/apinto-module"
)

type DynamicModuleDriver struct {
}

func NewDynamicModuleDriver() apinto_module.Driver {
	return &DynamicModuleDriver{}
}

func (c *DynamicModuleDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return &DynamicModulePlugin{
		define: define,
	}, nil
}

func (c *DynamicModulePlugin) GetPluginFrontend(moduleName string) string {
	return fmt.Sprintf("template/%s", moduleName)
}

func (c *DynamicModulePlugin) IsPluginVisible() bool {
	return true
}

func (c *DynamicModulePlugin) IsShowServer() bool {
	return false
}

func (c *DynamicModulePlugin) IsCanUninstall() bool {
	return true
}

func (c *DynamicModulePlugin) IsCanDisable() bool {
	return true
}

type DynamicModulePlugin struct {
	define interface{}
}

func (d *DynamicModulePlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewDynamicModule(name, d.define), nil
}

func (d *DynamicModulePlugin) CheckConfig(name string, config interface{}) error {
	return nil
}

type DynamicModule struct {
	name    string
	define  interface{}
	routers apinto_module.RoutersInfo
}

func (c *DynamicModule) Name() string {
	return c.name
}

func (c *DynamicModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *DynamicModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *DynamicModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewDynamicModule(name string, define interface{}) *DynamicModule {
	dm := &DynamicModule{name: name, define: define}
	dm.initRouter()
	return dm
}

func (c *DynamicModule) RoutersInfo() apinto_module.RoutersInfo {
	return c.routers
}
func (c *DynamicModule) initRouter() {
	dynamicController := newDynamicController(c.name, c.define)
	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/list", c.name),
			Handler:     "dynamic.list",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/info/:uuid", c.name),
			Handler:     "dynamic.info",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.info},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/render", c.name),
			Handler:     "dynamic.render",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.render},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/cluster/:uuid", c.name),
			Handler:     "dynamic.cluster_status",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.clusterStatus},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/status", c.name),
			Handler:     "dynamic.cluster_statuses",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.clusterStatusList},
		},
		{
			Method:      http.MethodDelete,
			Path:        fmt.Sprintf("/api/dynamic/%s/batch", c.name),
			Handler:     "dynamic.delete",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.batchDelete},
		}, {
			Method:      http.MethodPost,
			Path:        fmt.Sprintf("/api/dynamic/%s", c.name),
			Handler:     "dynamic.save",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.create},
		}, {
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/config/:uuid", c.name),
			Handler:     "dynamic.save",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.save},
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/online/:uuid", c.name),
			Handler:     "dynamic.online",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.online},
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/offline/:uuid", c.name),
			Handler:     "dynamic.offline",
			HandlerFunc: []apinto_module.HandlerFunc{dynamicController.offline},
		},
	}
}
