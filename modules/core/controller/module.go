package controller

import (
	"fmt"
	"net/http"

	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

var (
	_ apinto_module.Driver = (*Plugin)(nil)
	_ apinto_module.Plugin = (*Plugin)(nil)
	_ apinto_module.Module = (*Module)(nil)
)

type Plugin struct {
	middlewareHandler []apinto_module.MiddlewareHandler
	providers         apinto_module.IProviders
}

func NewCoreDriver() *Plugin {

	middlewareHandler := []apinto_module.MiddlewareHandler{
		{
			Name:    "namespace",
			Rule:    apinto_module.MiddlewareRule(apinto_module.RouterLabelApi),
			Handler: namespace_controller.MustNamespace,
		},
	}
	p := &Plugin{
		middlewareHandler: middlewareHandler,
	}

	bean.Autowired(&p.providers)
	return p
}

func (p *Plugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return p.NewModule(name), nil
}

func (p *Plugin) CheckConfig(name string, config interface{}) error {
	return nil
}

func (p *Plugin) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return p, nil
}

type Module struct {
	name              string
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
	return m, true
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
			name: result,
		},
	})

}
func (p *Plugin) NewModule(name string) *Module {

	routers := apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/provider/:name"),
			Handler:     "core.provider",
			HandlerFunc: []apinto_module.HandlerFunc{p.provider},
			Labels:      apinto_module.RouterLabelAssets,
		},
	}
	assets := staticFile("/assets", "dist/assets")
	routers = append(routers, assets...)
	aceBuilds := staticFile("/ace-builds", "dist/ace-builds")
	routers = append(routers, aceBuilds...)
	frontend := staticFile("/frontend", "dist")
	routers = append(routers, frontend...)

	routers = append(routers, favicon())
	routers = append(routers, indexRouter())
	routers = append(routers, commonStrategy()...)
	routers = append(routers, moduleRouters()...)
	systemRouter := newSystem()
	routers = append(routers, systemRouter.RoutersInfo()...)
	return &Module{
		name:              name,
		middlewareHandler: p.middlewareHandler,
		routers:           routers,
	}
}
