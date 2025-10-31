package audit_controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/apinto-dashboard/pm3/middleware"
	"github.com/eolinker/eosc/common/bean"
	"net/http"
)

type Driver struct {
	auditController   *auditLogController
	routers           apinto_module.RoutersInfo
	middlewareHandler []apinto_module.MiddlewareHandler
}

func (d *Driver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *Driver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return d.newModule(info.Id, info.Name), nil

}

func NewDriver() apinto_module.Driver {
	a := &auditLogController{}
	bean.Autowired(&a.auditLogService)

	return &Driver{
		auditController: a,
		routers: apinto_module.RoutersInfo{
			{
				Method: http.MethodGet,
				Path:   "/api/audit-logs",

				HandlerFunc: a.getLogs,
			},
			{
				Method: http.MethodGet,
				Path:   "/api/audit-log",

				HandlerFunc: a.getDetail,
			},
			{
				Method: http.MethodGet,
				Path:   "/api/audit-log/kinds",

				HandlerFunc: a.getTargets,
			},
		},
		middlewareHandler: []apinto_module.MiddlewareHandler{
			middleware.CreateF(a.Handler, func(api pm3.ApiInfo) bool {
				if api.Method == http.MethodGet {
					return false
				}
				switch api.Authority {
				case pm3.Private, pm3.Internal:
					return true
				}
				return false
			}),
		},
	}
}

func (d *Driver) newModule(id, name string) *Module {

	return &Module{
		ModuleTool: pm3.NewModuleTool(id, name),
		name:       name, middlewareHandler: d.middlewareHandler, routers: d.routers}
}

type Module struct {
	*pm3.ModuleTool

	name string

	routers           apinto_module.RoutersInfo
	middlewareHandler []apinto_module.MiddlewareHandler
}

func (m *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (m *Module) Middleware() []pm3.Middleware {
	return m.middlewareHandler
}

func (m *Module) Apis() []pm3.Api {
	m.InitAccess(m.routers)
	return m.routers
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}
