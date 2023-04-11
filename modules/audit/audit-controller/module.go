package audit_controller

import (
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"net/http"
)

type Driver struct {
	auditController   *auditLogController
	routers           apinto_module.RoutersInfo
	middlewareHandler []apinto_module.MiddlewareHandler
}

func NewDriver() *Driver {
	a := &auditLogController{}
	bean.Autowired(&a.auditLogService)

	return &Driver{
		auditController: a,
		routers: apinto_module.RoutersInfo{
			{
				Method:      http.MethodGet,
				Path:        "/api/audit-logs",
				Handler:     "audit.getLogs",
				HandlerFunc: []apinto_module.HandlerFunc{a.getLogs},
			},
			{
				Method:      http.MethodGet,
				Path:        "/api/audit-log",
				Handler:     "audit.getDetail",
				HandlerFunc: []apinto_module.HandlerFunc{a.getDetail},
			},
			{
				Method:      http.MethodGet,
				Path:        "/api/audit-log/kinds",
				Handler:     "audit.getTargets",
				HandlerFunc: []apinto_module.HandlerFunc{a.getTargets},
			},
		},
		middlewareHandler: []apinto_module.MiddlewareHandler{
			{
				Name:    "auditlog",
				Rule:    apinto_module.MiddlewareRules{{http.MethodPost}, {http.MethodPut}, {http.MethodGet, apinto_module.RouterTypeSensitive}},
				Handler: a.Handler,
			},
		},
	}
}

func (d *Driver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return d.newModule(name), nil
}

func (d *Driver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (d *Driver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return d, nil
}
func (d *Driver) newModule(name string) *Module {

	return &Module{name: name, middlewareHandler: d.middlewareHandler, routers: d.routers}
}

type Module struct {
	name string

	routers           apinto_module.RoutersInfo
	middlewareHandler []apinto_module.MiddlewareHandler
}

func (m *Module) MiddlewaresInfo() []apinto_module.MiddlewareHandler {
	return m.middlewareHandler
}

func (m *Module) RoutersInfo() apinto_module.RoutersInfo {
	return m.routers
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Routers() (apinto_module.Routers, bool) {
	return m, true
}

func (m *Module) Middleware() (apinto_module.Middleware, bool) {
	return m, true
}

func (m *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}
