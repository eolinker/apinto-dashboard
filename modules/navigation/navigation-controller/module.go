package navigation_controller

import (
	"net/http"

	"github.com/eolinker/apinto-module"
)

type NavigationDriver struct {
}

func NewNavigationPlugin() apinto_module.Driver {
	return &NavigationDriver{}
}

func (c *NavigationDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewNavigationModule(), nil
}

func (c *NavigationDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *NavigationDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type NavigationModule struct {
	isInit bool

	name    string
	routers apinto_module.RoutersInfo
}

func (c *NavigationModule) Name() string {
	return c.name
}

func (c *NavigationModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *NavigationModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *NavigationModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewNavigationModule() *NavigationModule {

	return &NavigationModule{}
}

func (c *NavigationModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}
func (c *NavigationModule) initRouter() {
	controller := newNavigationController()
	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/system/navigation",
			Handler:     "navigation.list",
			HandlerFunc: []apinto_module.HandlerFunc{controller.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/system/navigation/:uuid",
			Handler:     "navigation.info",
			HandlerFunc: []apinto_module.HandlerFunc{controller.info},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/system/navigation",
			Handler:     "navigation.create",
			HandlerFunc: []apinto_module.HandlerFunc{controller.add},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/system/navigation/:uuid",
			Handler:     "navigation.edit",
			HandlerFunc: []apinto_module.HandlerFunc{controller.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/system/navigation/:uuid",
			Handler:     "navigation.delete",
			HandlerFunc: []apinto_module.HandlerFunc{controller.delete},
		}, {
			Method:      http.MethodPut,
			Path:        "/api/system/navigation",
			Handler:     "navigation.sort",
			HandlerFunc: []apinto_module.HandlerFunc{controller.sort},
		},
	}
}
