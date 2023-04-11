package main

import (
	apiController "github.com/eolinker/apinto-dashboard/modules/api/controller"
	auditController "github.com/eolinker/apinto-dashboard/modules/audit/audit-controller"
	clusterController "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	coreController "github.com/eolinker/apinto-dashboard/modules/core/controller"
	discovery_controller "github.com/eolinker/apinto-dashboard/modules/discovery/discovery-controller"
	navigationController "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-controller"
	open_api_controller "github.com/eolinker/apinto-dashboard/modules/openapi/open-api-controller"
	apintoModule "github.com/eolinker/apinto-module"
)

func init() {

	apintoModule.Register("core.apinto.com", coreController.NewCoreDriver())
	apintoModule.Register("api.apinto.com", apiController.NewPluginDriver())
	apintoModule.Register("cluster.apinto.com", clusterController.NewClusterPlugin())
	apintoModule.Register("navigation.apinto.com", navigationController.NewNavigationPlugin())
	apintoModule.Register("audit.apinto.com", auditController.NewDriver())
	apintoModule.Register("discovery.apinto.com", discovery_controller.NewPluginDriver())
	apintoModule.Register("open_api.apinto.com", open_api_controller.NewPluginDriver())
}
