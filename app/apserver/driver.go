package main

import (
	audit_controller "github.com/eolinker/apinto-dashboard/modules/audit/audit-controller"
	cluster_controller "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	"github.com/eolinker/apinto-dashboard/modules/core/controller"
	navigation_controller "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-controller"
	apinto_module "github.com/eolinker/apinto-module"
)

func init() {

	apinto_module.Register("core.apinto.com", controller.NewCoreDriver())
	apinto_module.Register("cluster.apinto.com", cluster_controller.NewClusterPlugin())
	apinto_module.Register("navigation.apinto.com", navigation_controller.NewNavigationPlugin())
	apinto_module.Register("audit.apinto.com", audit_controller.NewDriver())
}
