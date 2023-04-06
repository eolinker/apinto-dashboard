package main

import (
	cluster_controller "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-controller"
	apinto_module "github.com/eolinker/apinto-module"
)

func init() {
	apinto_module.Register("cluster.apinto.com", cluster_controller.NewClusterPlugin())
}
