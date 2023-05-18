package cluster_model

import (
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
)

type ClusterCertificate struct {
	*cluster_entry.ClusterCertificate
	ClusterName  string
	OperatorName string
}

type Certificate struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	DnsName   []string `json:"dns_name"`
	Key       string   `json:"key"`
	Pem       string   `json:"pem"`
	ValidTime string   `json:"valid_time"`
}
