package controller

import (
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

type ModulePluginDriver struct {
}

func NewModulePlugin() apinto_module.Driver {
	return &ModulePluginDriver{}
}

func (c *ModulePluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewModulePluginModule(name), nil
}

func (c *ModulePluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *ModulePluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type ModulePluginModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *ModulePluginModule) Name() string {
	return c.name
}

func (c *ModulePluginModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *ModulePluginModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *ModulePluginModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewModulePluginModule(name string) *ModulePluginModule {

	return &ModulePluginModule{name: name}
}

func (c *ModulePluginModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *ModulePluginModule) initRouter() {
	mPluginController := newModulePluginController()
	controllerPluginFront := newPluginFrontController()
	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/installed",
			Handler:     "modulePlugin.plugins",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.plugins},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/info",
			Handler:     "modulePlugin.getPluginInfo",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.getPluginInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/groups/enum",
			Handler:     "modulePlugin.getGroupsEnum",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.getGroupsEnum},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/plugin/enable",
			Handler:     "modulePlugin.getEnableInfo",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.getEnableInfo},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/plugin/install",
			Handler:     "modulePlugin.install",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.install},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/plugin/enable",
			Handler:     "modulePlugin.enable",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.enable},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/plugin/disable",
			Handler:     "modulePlugin.disable",
			HandlerFunc: []apinto_module.HandlerFunc{mPluginController.disable},
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/icon/:id/:file",
			Handler:     "pluginFront.getPluginIcon",
			HandlerFunc: []apinto_module.HandlerFunc{controllerPluginFront.checkPluginID, controllerPluginFront.setIConName, controllerPluginFront.getPluginInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/icon/:id",
			Handler:     "pluginFront.getPluginIconDefault",
			HandlerFunc: []apinto_module.HandlerFunc{controllerPluginFront.checkPluginID, controllerPluginFront.setIConName, controllerPluginFront.getPluginInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/md/:id/:file",
			Handler:     "pluginFront.getPluginMD",
			HandlerFunc: []apinto_module.HandlerFunc{controllerPluginFront.checkPluginID, controllerPluginFront.setMDName, controllerPluginFront.getPluginInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/md/:id",
			Handler:     "pluginFront.getPluginMDDefault",
			HandlerFunc: []apinto_module.HandlerFunc{controllerPluginFront.checkPluginID, controllerPluginFront.setMDName, controllerPluginFront.getPluginInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/info/:id/resources/*filepath",
			Handler:     "pluginFront.getPluginResources",
			HandlerFunc: []apinto_module.HandlerFunc{controllerPluginFront.checkPluginID, controllerPluginFront.getPluginResources},
		},
	}
}
