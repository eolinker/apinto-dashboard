package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
)

type ClusterCertificate struct {
	*cluster_entry.ClusterCertificate
	ClusterName  string
	OperatorName string
}
