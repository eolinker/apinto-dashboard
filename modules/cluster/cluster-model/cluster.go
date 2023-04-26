package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
)

type Cluster struct {
	*cluster_entry.Cluster
	Status int //1正常 2部分正常 3异常
}

type ClusterSimple struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
