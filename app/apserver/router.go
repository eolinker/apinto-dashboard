package main

import (
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/application/application-controller"
	"github.com/eolinker/apinto-dashboard/modules/base/env-controller"
	"github.com/eolinker/apinto-dashboard/modules/base/random-controller"
	cluster_controller2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discovery-controller"
	"github.com/eolinker/apinto-dashboard/modules/group/group-controller"
	"github.com/eolinker/apinto-dashboard/modules/openapi/open-api-controller"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-controller"
	plugin_controller "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-controller"
	plugin_template_controller "github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-controller"
	strategy_controller2 "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-controller"
	user_controller "github.com/eolinker/apinto-dashboard/modules/user/user-controller"
	variable_controller2 "github.com/eolinker/apinto-dashboard/modules/variable/variable-controller"

	"net/http"

	apiController "github.com/eolinker/apinto-dashboard/modules/api/controller"
	upstream_controller "github.com/eolinker/apinto-dashboard/modules/upstream/controller"
	"github.com/gin-gonic/gin"
)

//
//var mustNamespaceExclude = []string{"/api/random/:template/id", "/api/enum/envs", "/api/application/drivers",
//	"/api/access", "/api/my/modules", "/api/my/profile", "/api/my/password", "/api/roles", "/api/role", "/api/user/password-reset",
//	"/api/role/options", "/api/role", "/api/role/batch-update", "/api/role/batch-remove", "/api/user/delete", "/api/user/profile", "/api/user/list"}
//
//var verifyTokenExclude = []string{"/_system/activation", "/_system/mac", "/_system/run-info"}

func registerRouter(engine *gin.Engine) {

	engine.Use()

	routes := engine.Group("/api")

	cluster_controller2.RegisterClusterRouter(routes)
	cluster_controller2.RegisterClusterCertificateRouter(routes)
	variable_controller2.RegisterClusterVariableRouter(routes)
	cluster_controller2.RegisterClusterNodeRouter(routes)
	cluster_controller2.RegisterClusterConfigRouter(routes)
	variable_controller2.RegisterVariablesRouter(routes)
	upstream_controller.RegisterServiceRouter(routes)
	discovery_controller.RegisterDiscoveryRouter(routes)
	application_controller.RegisterApplicationRouter(routes)
	group_controller.RegisterCommonGroupRouter(routes)

	//middleware_controller.RegisterMiddlewareGroupRouter(routes)
	env_controller.RegisterEnumRouter(routes)
	apiController.RegisterAPIRouter(routes)
	random_controller.RegisterRandomRouter(routes)

	strategy_controller2.RegisterStrategyCommonRouter(routes)
	strategy_controller2.RegisterStrategyTrafficRouter(routes)
	strategy_controller2.RegisterStrategyCacheRouter(routes)
	strategy_controller2.RegisterStrategyGreyRouter(routes)
	strategy_controller2.RegisterStrategyVisitRouter(routes)
	strategy_controller2.RegisterStrategyFuseRouter(routes)
	plugin_controller.RegisterPluginRouter(routes)
	plugin_controller.RegisterPluginClusterRouter(routes)
	plugin_template_controller.RegisterPluginTemplateRouter(routes)
	open_app_controller.RegisterExternalApplicationRouter(routes)

	user_controller.RegisterUserRouter(routes)

	controller.EmbedFrontend(engine)

	openAPIRoutes := engine.Group("/api2")
	open_api_controller.RegisterApiOpenAPIRouter(openAPIRoutes) //api管理导入的OpenAPI

	engine.Handle(http.MethodGet, "/_system/profile", version.Handler)
}
