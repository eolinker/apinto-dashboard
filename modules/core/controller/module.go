package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/frontend"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/apinto-dashboard/pm3/middleware"
	"net/http"
	"strings"

	"github.com/eolinker/eosc/common/bean"

	apinto_module "github.com/eolinker/apinto-dashboard/module"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	notice_controller "github.com/eolinker/apinto-dashboard/modules/notice/controller"
	"github.com/gin-gonic/gin"
)

var (
	_ apinto_module.Module = (*Module)(nil)
)

type Module struct {
	*pm3.ModuleTool

	middlewareHandler []pm3.Middleware
	frontendAssets    []pm3.FrontendAsset
	apis              []pm3.Api
	providers         apinto_module.IProviders
}

func (m *Module) Frontend() []pm3.FrontendAsset {
	return m.frontendAssets
}

func (m *Module) Apis() []pm3.Api {
	return m.apis
}

func (m *Module) Middleware() []pm3.Middleware {
	return m.middlewareHandler
}

func (m *Module) Name() string {
	return "core"
}

func (m *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (m *Module) provider(context *gin.Context) {
	skill := context.Param("skill")
	namespaceID := namespace_controller.GetNamespaceId(context)
	provider, ok := m.providers.Provider(skill)
	if !ok {
		context.JSON(200, struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}{
			"200", fmt.Sprintf("not support data for %s", skill),
		})
		return
	}
	cargos := provider.Provide(namespaceID)
	result := make([]*pm3.CargoItem, 0, len(cargos))
	for _, c := range cargos {
		result = append(result, c.Export())
	}
	context.JSON(200, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			skill: result,
		},
	})

}

func (m *Module) enum(context *gin.Context) {
	skill := context.Param("skill")
	namespaceID := namespace_controller.GetNamespaceId(context)
	provider, ok := m.providers.Provider(skill)
	if !ok {
		context.JSON(200, struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}{
			"200", fmt.Sprintf("not support data for %s", skill),
		})
		return
	}
	cargos := provider.Provide(namespaceID)
	result := make([]*pm3.CargoItem, 0, len(cargos))
	for _, c := range cargos {
		export := c.Export()
		index := strings.LastIndex(export.Value, "@")
		if index > 0 {
			export.Value = export.Value[:index]
		}
		result = append(result, export)
	}
	context.JSON(200, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			skill: result,
		},
	})
}

func (m *Module) status(context *gin.Context) {
	key := context.Query("name")
	if key == "" {
		context.JSON(200, struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}{
			"200", "empty name",
		})
		return
	}
	cluster := context.Query("cluster")
	if key == "" {
		context.JSON(200, struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
		}{
			"200", "empty cluster",
		})
		return
	}

	namespaceID := namespace_controller.GetNamespaceId(context)

	status, _ := m.providers.Status(key, namespaceID, cluster)

	context.JSON(200, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"status": status,
		},
	})

}
func NewModule(id, name string) *Module {
	middlewareHandler := []apinto_module.MiddlewareHandler{
		middleware.Create(namespace_controller.MustNamespace, middleware.IsApi),
	}

	m := &Module{
		ModuleTool:        pm3.NewModuleTool(id, name),
		middlewareHandler: middlewareHandler,
		frontendAssets:    frontend.Frontends(),
	}
	routers := make([]pm3.Api, 0, 25)
	routers = apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/provider/:skill"),
			HandlerFunc: m.provider,
			Authority:   pm3.Public,
		},
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/common/status"),

			HandlerFunc: m.status,
			Authority:   pm3.Public,
		}, {
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/common/enum/:skill"),

			HandlerFunc: m.enum,
			Authority:   pm3.Public,
		}, {
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/report"),
			HandlerFunc: m.reportStatus,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodPost,
			Path:        fmt.Sprintf("/api/common/report"),
			HandlerFunc: m.updateReport,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/common/report"),
			HandlerFunc: m.updateReport,
			Authority:   pm3.Public,
		},
	}

	routers = append(routers, commonStrategy()...)
	systemRouter := systemRouters()
	routers = append(routers, systemRouter...)
	routers = append(routers, envEnumRouters()...)
	routers = append(routers, notice_controller.InitRouter()...)
	routers = append(routers, randomRouters()...)
	routers = append(routers, userApis()...)
	m.apis = routers
	m.InitAccess(m.apis)
	bean.Autowired(&m.providers)
	return m
}
