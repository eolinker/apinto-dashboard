package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"net/http"
)

type StrategyTrafficDriver struct {
}

func NewStrategyTraffic() apinto_module.Driver {
	return &StrategyTrafficDriver{}
}

func (c *StrategyTrafficDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewStrategyTrafficModule(name), nil
}

func (c *StrategyTrafficDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *StrategyTrafficDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *StrategyTrafficDriver) GetPluginFrontend(moduleName string) string {
	return "serv-governance/traffic"
}

func (c *StrategyTrafficDriver) IsPluginVisible() bool {
	return true
}

func (c *StrategyTrafficDriver) IsShowServer() bool {
	return false
}

func (c *StrategyTrafficDriver) IsCanUninstall() bool {
	return false
}

func (c *StrategyTrafficDriver) IsCanDisable() bool {
	return true
}

type StrategyTrafficModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *StrategyTrafficModule) Name() string {
	return c.name
}

func (c *StrategyTrafficModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *StrategyTrafficModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *StrategyTrafficModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewStrategyTrafficModule(name string) *StrategyTrafficModule {

	return &StrategyTrafficModule{name: name}
}

func (c *StrategyTrafficModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *StrategyTrafficModule) initRouter() {
	strategyTrafficController := newStrategyTrafficController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategies/traffic",
			Handler:     "strategy-traffic.list",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.list},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/traffic",
			Handler:     "strategy-traffic.get",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/traffic",
			Handler:     "strategy-traffic.create",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/strategy/traffic",
			Handler:     "strategy-traffic.update",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/strategy/traffic",
			Handler:     "strategy-traffic.del",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.del},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/traffic/restore",
			Handler:     "strategy-traffic.restore",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.restore},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/api/strategy/traffic/stop",
			Handler:     "strategy-traffic.updateStop",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.updateStop},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/traffic/to-publishs",
			Handler:     "strategy-traffic.toPublish",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.toPublish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/traffic/publish",
			Handler:     "strategy-traffic.publish",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, strategyTrafficController.publish},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/strategy/traffic/priority",
			Handler:     "strategy-traffic.changePriority",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypeEdit.Handler, strategyTrafficController.changePriority},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/traffic/publish-history",
			Handler:     "strategy-traffic.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{strategyTrafficController.publishHistory},
		},
	}
}
