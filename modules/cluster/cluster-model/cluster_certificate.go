package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
)

type ClusterCertificate struct {
	*cluster_entry.ClusterCertificate
	ClusterName  string
	OperatorName string
}
