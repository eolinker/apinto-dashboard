package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
)

type Cluster struct {
	*cluster_entry.Cluster
	Status int //1正常 2部分正常 3异常
}
