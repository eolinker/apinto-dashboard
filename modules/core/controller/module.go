package controller

import (
	"fmt"
	"net/http"

	"github.com/eolinker/eosc/common/bean"

	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	notice_controller "github.com/eolinker/apinto-dashboard/modules/notice/controller"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
)

var (
	_ apinto_module.Module = (*Module)(nil)
)

type Module struct {
	middlewareHandler []apinto_module.MiddlewareHandler
	routers           apinto_module.RoutersInfo

	providers apinto_module.IProviders
}

func (m *Module) RoutersInfo() apinto_module.RoutersInfo {
	return m.routers
}

func (m *Module) MiddlewaresInfo() []apinto_module.MiddlewareHandler {
	return m.middlewareHandler
}

func (m *Module) Name() string {
	return "core"
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
	result := make([]*apinto_module.CargoItem, 0, len(cargos))
	for _, c := range cargos {
		result = append(result, c.Export())
	}
	context.JSON(200, map[string]interface{}{
		"code": "00000",
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

	status := m.providers.Status(key, namespaceID, cluster)

	context.JSON(200, map[string]interface{}{
		"code": "00000",
		"data": map[string]interface{}{
			"status": status,
		},
	})

}
func NewModule() *Module {
	userController := newUserController()
	middlewareHandler := []apinto_module.MiddlewareHandler{
		{
			Name:    "namespace",
			Rule:    apinto_module.MiddlewareRule(apinto_module.RouterLabelApi),
			Handler: namespace_controller.MustNamespace,
		},
		{
			Name:    "userID",
			Rule:    apinto_module.MiddlewareRule(apinto_module.RouterLabelApi),
			Handler: userController.SetUser,
		}, {
			Name:    "login-api",
			Rule:    apinto_module.MiddlewareRule(apinto_module.RouterLabelApi),
			Handler: userController.LoginCheckApi,
		}, {
			Name:    "login-api",
			Rule:    apinto_module.MiddlewareRule(apinto_module.RouterLabelApi),
			Handler: userController.LoginCheckApi,
		},
	}
	m := &Module{

		middlewareHandler: middlewareHandler,
	}
	routers := apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/provider/:skill"),
			Handler:     "core.provider",
			HandlerFunc: []apinto_module.HandlerFunc{m.provider},
			Labels:      apinto_module.RouterLabelAssets,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/common/status"),
			Handler:     "core.provider",
			HandlerFunc: []apinto_module.HandlerFunc{m.status},
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
	routers = append(routers, envEnumRouters()...)
	routers = append(routers, notice_controller.InitRouter()...)
	routers = append(routers, userRouters()...)
	routers = append(routers, randomRouters()...)
	m.routers = routers

	bean.Autowired(&m.providers)
	return m
}
