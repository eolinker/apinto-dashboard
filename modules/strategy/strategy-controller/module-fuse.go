package strategy_controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-module"
	"net/http"
)

type StrategyFuseDriver struct {
}

func NewStrategyFuse() apinto_module.Driver {
	return &StrategyFuseDriver{}
}

func (c *StrategyFuseDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewStrategyFuseModule(name), nil
}

func (c *StrategyFuseDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *StrategyFuseDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *StrategyFuseDriver) GetPluginFrontend(moduleName string) string {
	return "serv-governance/fuse"
}

func (c *StrategyFuseDriver) IsPluginVisible() bool {
	return true
}

func (c *StrategyFuseDriver) IsShowServer() bool {
	return false
}

func (c *StrategyFuseDriver) IsCanUninstall() bool {
	return false
}

func (c *StrategyFuseDriver) IsCanDisable() bool {
	return true
}

type StrategyFuseModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyFuseModule) Name() string {
	return c.name
}

func (c *StrategyFuseModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *StrategyFuseModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *StrategyFuseModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewStrategyFuseModule(name string) *StrategyFuseModule {

	return &StrategyFuseModule{name: name}
}

func (c *StrategyFuseModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *StrategyFuseModule) initRouter() {
	strategyFuseController := newStrategyFuseController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategies/fuse",
			Handler:     "strategy-fuse.list",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/fuse",
			Handler:     "strategy-fuse.get",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/fuse",
			Handler:     "strategy-fuse.create",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/strategy/fuse",
			Handler:     "strategy-fuse.update",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/strategy/fuse",
			Handler:     "strategy-fuse.del",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.del},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/fuse/restore",
			Handler:     "strategy-fuse.restore",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.restore},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/fuse/stop",
			Handler:     "strategy-fuse.updateStop",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.updateStop},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/fuse/to-publishs",
			Handler:     "strategy-fuse.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.toPublish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/fuse/publish",
			Handler:     "strategy-fuse.publish",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, strategyFuseController.publish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/fuse/priority",
			Handler:     "strategy-fuse.changePriority",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypeEdit.Handler, strategyFuseController.changePriority},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/fuse/publish-history",
			Handler:     "strategy-fuse.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{strategyFuseController.publishHistory},
		},
	}
}
