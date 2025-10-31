package main

import (
	apintoModule "github.com/eolinker/apinto-dashboard/module"
	apiController "github.com/eolinker/apinto-dashboard/modules/api/controller"
	appController "github.com/eolinker/apinto-dashboard/modules/application/application-controller"
	auditController "github.com/eolinker/apinto-dashboard/modules/audit/audit-controller"
	clusterController "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	_ "github.com/eolinker/apinto-dashboard/modules/core/controller"
	dynamic_controller "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-controller"
	email_controller "github.com/eolinker/apinto-dashboard/modules/email/controller"
	logModule "github.com/eolinker/apinto-dashboard/modules/log/module"
	module_plugin_controller "github.com/eolinker/apinto-dashboard/modules/mpm3/controller"
	open_api_controller "github.com/eolinker/apinto-dashboard/modules/openapi/open-api-controller"
	open_app_controller "github.com/eolinker/apinto-dashboard/modules/openapp/open-app-controller"
	plugin_controller "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-controller"
	plugin_template_controller "github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-controller"
	strategy_controller "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-controller"
	variable_controller "github.com/eolinker/apinto-dashboard/modules/variable/variable-controller"
	webhook_controller "github.com/eolinker/apinto-dashboard/modules/webhook/controller"
	"github.com/eolinker/apinto-dashboard/plugin/local"
	"github.com/eolinker/apinto-dashboard/plugin/remote"
)

func init() {

	apintoModule.Register("api.apinto.com", apiController.NewPluginDriver())
	apintoModule.Register("application.apinto.com", appController.NewPluginDriver())
	apintoModule.Register("audit.apinto.com", auditController.NewDriver())
	apintoModule.Register("cluster.apinto.com", clusterController.NewClusterPlugin())
	apintoModule.Register("ext_app.apinto.com", open_app_controller.NewPluginDriver())
	apintoModule.Register("module_plugin.apinto.com", module_plugin_controller.NewDriver())
	apintoModule.Register("open_api.apinto.com", open_api_controller.NewPluginDriver())
	apintoModule.Register("plugin.apinto.com", plugin_controller.NewPluginDriver())
	apintoModule.Register("plugin_template.apinto.com", plugin_template_controller.NewPluginTemplateDriver())
	apintoModule.Register("strategy-cache.apinto.com", strategy_controller.NewStrategyCache())
	apintoModule.Register("strategy-fuse.apinto.com", strategy_controller.NewStrategyFuse())
	apintoModule.Register("strategy-grey.apinto.com", strategy_controller.NewStrategyGrey())
	apintoModule.Register("strategy-traffic.apinto.com", strategy_controller.NewStrategyTraffic())
	apintoModule.Register("strategy-visit.apinto.com", strategy_controller.NewStrategyVisit())
	apintoModule.Register("strategy-data_mask.apinto.com", strategy_controller.NewStrategyDataMask())
	apintoModule.Register("variable.apinto.com", variable_controller.NewVariableDriver())

	apintoModule.Register("local", local.NewDriver())
	apintoModule.Register("remote", remote.NewDriver())

	apintoModule.Register("email.apinto.com", email_controller.NewEmailDriver())
	apintoModule.Register("webhook.apinto.com", webhook_controller.NewWebhookDriver())
	apintoModule.Register("log.apinto.com", logModule.NewDriver())

	apintoModule.Register("dynamic.apinto.com", dynamic_controller.NewDynamicModuleDriver())
	apintoModule.Register("upstream.apinto.com", dynamic_controller.NewDynamicModuleDriver())
	apintoModule.Register("discovery.apinto.com", dynamic_controller.NewDynamicModuleDriver())
	//apintoModule.Register("application.apinto.com", dynamic_controller.NewDynamicModuleDriver(true, false, false, false))

}
