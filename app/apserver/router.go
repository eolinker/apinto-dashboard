package main

import (
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/application-controller"
	"github.com/eolinker/apinto-dashboard/controller/audit-controller"
	"github.com/eolinker/apinto-dashboard/controller/bussiness-controller"
	"github.com/eolinker/apinto-dashboard/controller/cluster-controller"
	"github.com/eolinker/apinto-dashboard/controller/common/namespace-controller"
	"github.com/eolinker/apinto-dashboard/controller/discovery-controller"
	"github.com/eolinker/apinto-dashboard/controller/env-controller"
	"github.com/eolinker/apinto-dashboard/controller/frontend-controller"
	"github.com/eolinker/apinto-dashboard/controller/group-controller"
	"github.com/eolinker/apinto-dashboard/controller/monitor-controller"
	"github.com/eolinker/apinto-dashboard/controller/open-api-controller"
	"github.com/eolinker/apinto-dashboard/controller/open-app-controller"
	"github.com/eolinker/apinto-dashboard/controller/random-controller"
	"github.com/eolinker/apinto-dashboard/controller/strategy-controller"
	user_center "github.com/eolinker/apinto-dashboard/controller/user-center"
	"github.com/eolinker/apinto-dashboard/controller/user-controller"
	"github.com/eolinker/apinto-dashboard/controller/variable-controller"
	"github.com/eolinker/apinto-dashboard/filter"
	apiController "github.com/eolinker/apinto-dashboard/modules/api/controller"
	upstream_controller "github.com/eolinker/apinto-dashboard/modules/upstream/controller"
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

	engine.Use(controller.Logger, controller.Recovery, verifyTokenBuilder.Build(controller.VerifyToken, controller.VerifyAuth), namespaceBuilder.Build(namespace_controller.MustNamespace))

	routes := engine.Group("/api")
	user_center.RegisterUserCenterProxyRouter(engine)
	cluster_controller.RegisterClusterRouter(routes)
	cluster_controller.RegisterClusterCertificateRouter(routes)
	variable_controller.RegisterClusterVariableRouter(routes)
	cluster_controller.RegisterClusterNodeRouter(routes)
	cluster_controller.RegisterClusterConfigRouter(routes)
	variable_controller.RegisterVariablesRouter(routes)
	upstream_controller.RegisterServiceRouter(routes)
	discovery_controller.RegisterDiscoveryRouter(routes)
	application_controller.RegisterApplicationRouter(routes)
	group_controller.RegisterCommonGroupRouter(routes)
	env_controller.RegisterEnumRouter(routes)
	apiController.RegisterAPIRouter(routes)
	random_controller.RegisterRandomRouter(routes)

	strategy_controller.RegisterStrategyCommonRouter(routes)
	strategy_controller.RegisterStrategyTrafficRouter(routes)
	strategy_controller.RegisterStrategyCacheRouter(routes)
	strategy_controller.RegisterStrategyGreyRouter(routes)
	strategy_controller.RegisterStrategyVisitRouter(routes)
	strategy_controller.RegisterStrategyFuseRouter(routes)
	audit_controller.RegisterAuditLogRouter(routes)
	open_app_controller.RegisterExternalApplicationRouter(routes)
	monitor_controller.RegisterMonitorRouter(routes)

	user_controller.RegisterUserRouter(routes)

	monitor_controller.RegisterWarnRouter(routes)

	//注册商业授权路由
	bussiness_controller.RegisterBussinessAuthRouter(engine)

	frontend_controller.EmbedFrontend(engine)

	openAPIRoutes := engine.Group("/api2")
	open_api_controller.RegisterApiOpenAPIRouter(openAPIRoutes) //api管理导入的OpenAPI

	engine.Handle(http.MethodGet, "/_system/profile", version.Handler)
}
