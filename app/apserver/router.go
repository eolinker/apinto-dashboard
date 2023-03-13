package main

import (
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/controller"
	user_center "github.com/eolinker/apinto-dashboard/controller/user-center"
	"github.com/eolinker/apinto-dashboard/filter"
	"github.com/gin-gonic/gin"
	"net/http"
)

var mustNamespaceExclude = []string{"/api/random/:template/id", "/api/enum/envs", "/api/application/drivers",
	"/api/access", "/api/my/modules", "/api/my/profile", "/api/my/password", "/api/roles", "/api/role", "/api/user/password-reset",
	"/api/role/options", "/api/role", "/api/role/batch-update", "/api/role/batch-remove", "/api/user/delete", "/api/user/profile", "/api/user/list"}

var verifyTokenExclude = []string{"/_system/activation", "/_system/mac", "/_system/run-info"}

func registerRouter(engine *gin.Engine) {
	namespaceBuilder := filter.NewBuilder().Root("/api/").Exclude(mustNamespaceExclude...)
	verifyTokenBuilder := filter.NewBuilder().Root("/api/", "/_system/").Exclude(verifyTokenExclude...)

	engine.Use(controller.Logger, controller.Recovery, verifyTokenBuilder.Build(controller.VerifyToken, controller.VerifyAuth), namespaceBuilder.Build(controller.MustNamespace))

	routes := engine.Group("/api")
	user_center.RegisterUserCenterProxyRouter(engine)
	controller.RegisterClusterRouter(routes)
	controller.RegisterClusterCertificateRouter(routes)
	controller.RegisterClusterVariableRouter(routes)
	controller.RegisterClusterNodeRouter(routes)
	controller.RegisterClusterConfigRouter(routes)
	controller.RegisterVariablesRouter(routes)
	controller.RegisterServiceRouter(routes)
	controller.RegisterDiscoveryRouter(routes)
	controller.RegisterApplicationRouter(routes)
	controller.RegisterCommonGroupRouter(routes)
	controller.RegisterEnumRouter(routes)
	controller.RegisterAPIRouter(routes)
	controller.RegisterRandomRouter(routes)

	controller.RegisterStrategyCommonRouter(routes)
	controller.RegisterStrategyTrafficRouter(routes)
	controller.RegisterStrategyCacheRouter(routes)
	controller.RegisterStrategyGreyRouter(routes)
	controller.RegisterStrategyVisitRouter(routes)
	controller.RegisterStrategyFuseRouter(routes)
	controller.RegisterAuditLogRouter(routes)
	controller.RegisterExternalApplicationRouter(routes)
	controller.RegisterMonitorRouter(routes)

	controller.RegisterUserRouter(routes)

	controller.RegisterWarnRouter(routes)

	//注册商业授权路由
	controller.RegisterBussinessAuthRouter(engine)

	controller.EmbedFrontend(engine)

	openAPIRoutes := engine.Group("/api2")
	controller.RegisterApiOpenAPIRouter(openAPIRoutes) //api管理导入的OpenAPI

	engine.Handle(http.MethodGet, "/_system/profile", version.Handler)
}
