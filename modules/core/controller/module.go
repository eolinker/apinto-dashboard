package controller

import (
	"fmt"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/core"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	_ apinto_module.Driver = (*Plugin)(nil)
	_ apinto_module.Plugin = (*Plugin)(nil)
	_ apinto_module.Module = (*Module)(nil)
)

type Plugin struct {
	middlewareHandler []apinto_module.MiddlewareHandler
	providers         core.IProviders
}

func NewCoreDriver() *Plugin {

	middlewareHandler := []apinto_module.MiddlewareHandler{
		{
			Name:    "namespace",
			Handler: namespace_controller.MustNamespace,
		},
	}
	p := &Plugin{
		middlewareHandler: middlewareHandler,
	}

	bean.Autowired(&p.providers)
	return p
}

func (p *Plugin) CreateModule(name string, apiPrefix string, config interface{}) (apinto_module.Module, error) {
	return p.NewModule(name, apiPrefix), nil
}

func (p *Plugin) CheckConfig(name string, apiPrefix string, config interface{}) error {
	return nil
}

func (p *Plugin) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return p, nil
}

type Module struct {
	name              string
	apiPrefix         string
	middlewareHandler []apinto_module.MiddlewareHandler
	routers           apinto_module.RoutersInfo
}

func (m *Module) RoutersInfo() apinto_module.RoutersInfo {
	return m.routers
}

func (m *Module) MiddlewaresInfo() []apinto_module.MiddlewareHandler {
	return m.middlewareHandler
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Routers() (apinto_module.Routers, bool) {
	return m, false
}

func (m *Module) Middleware() (apinto_module.Middleware, bool) {
	return m, true
}

func (m *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}
func (p *Plugin) provider(context *gin.Context) {
	name := context.Param("name")

	provider, ok := p.providers.Provider(name)
	if !ok {
		context.JSON(200, struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}{
			"200", fmt.Sprintf("not support data for %s", name),
		})
	}
	cargos := provider.Provide()
	result := make([]*apinto_module.CargoItem, 0, len(cargos))
	for _, c := range cargos {
		result = append(result, c.Export())
	}
	context.JSON(200, map[string]interface{}{
		"code": "00000",
		"data": map[string]interface{}{
			"cargos": result,
		},
	})

}
func (p *Plugin) NewModule(name, apiPrefix string) *Module {

	routers := apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/provider/:name"),
			Handler:     "core.provider",
			HandlerFunc: []apinto_module.HandlerFunc{p.provider},
		},
	}

	return &Module{
		name:              name,
		apiPrefix:         apiPrefix,
		middlewareHandler: p.middlewareHandler,
		routers:           routers,
	}
}
