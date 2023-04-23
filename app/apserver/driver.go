package main

import (
	apiController "github.com/eolinker/apinto-dashboard/modules/api/controller"
	application_controller "github.com/eolinker/apinto-dashboard/modules/application/application-controller"
	auditController "github.com/eolinker/apinto-dashboard/modules/audit/audit-controller"
	clusterController "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	coreController "github.com/eolinker/apinto-dashboard/modules/core/controller"
	discovery_controller "github.com/eolinker/apinto-dashboard/modules/discovery/discovery-controller"
	module_plugin_controller "github.com/eolinker/apinto-dashboard/modules/module-plugin/controller"
	open_api_controller "github.com/eolinker/apinto-dashboard/modules/openapi/open-api-controller"
	open_app_controller "github.com/eolinker/apinto-dashboard/modules/openapp/open-app-controller"
	plugin_controller "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-controller"
	plugin_template_controller "github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-controller"
	strategy_controller "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-controller"
	upstream_controller "github.com/eolinker/apinto-dashboard/modules/upstream/controller"
	variable_controller "github.com/eolinker/apinto-dashboard/modules/variable/variable-controller"
	"github.com/eolinker/apinto-dashboard/plugin/local"
	apintoModule "github.com/eolinker/apinto-module"
)

func init() {

	apintoModule.Register("api.apinto.com", apiController.NewPluginDriver())
	apintoModule.Register("application.apinto.com", application_controller.NewPluginDriver())
	apintoModule.Register("audit.apinto.com", auditController.NewDriver())
	apintoModule.Register("cluster.apinto.com", clusterController.NewClusterPlugin())
	apintoModule.Register("core.apinto.com", coreController.NewCoreDriver())
	apintoModule.Register("discovery.apinto.com", discovery_controller.NewPluginDriver())
	apintoModule.Register("ext_app.apinto.com", open_app_controller.NewPluginDriver())
	apintoModule.Register("module_plugin.apinto.com", module_plugin_controller.NewModulePlugin())
	apintoModule.Register("open_api.apinto.com", open_api_controller.NewPluginDriver())
	apintoModule.Register("plugin.apinto.com", plugin_controller.NewPluginDriver())
	apintoModule.Register("plugin_template.apinto.com", plugin_template_controller.NewPluginTemplateDriver())
	apintoModule.Register("strategy-cache.apinto.com", strategy_controller.NewStrategyCache())
	apintoModule.Register("strategy-fuse.apinto.com", strategy_controller.NewStrategyFuse())
	apintoModule.Register("strategy-grey.apinto.com", strategy_controller.NewStrategyGrey())
	apintoModule.Register("strategy-traffic.apinto.com", strategy_controller.NewStrategyTraffic())
	apintoModule.Register("strategy-visit.apinto.com", strategy_controller.NewStrategyVisit())
	apintoModule.Register("upstream.apinto.com", upstream_controller.NewUpstreamDriver())
	apintoModule.Register("variable.apinto.com", variable_controller.NewVariableDriver())
	apintoModule.Register("local", local.NewDriver())
}
