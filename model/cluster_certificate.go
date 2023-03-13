package model

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type ClusterCertificate struct {
	*entry.ClusterCertificate
	ClusterName  string
	OperatorName string
}
