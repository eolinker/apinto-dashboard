package dto

import "github.com/eolinker/apinto-dashboard/enum"

type ClusterNode struct {
	Name        string                 `json:"name"`
	ServiceAddr string                 `json:"service_addr"`
	AdminAddr   string                 `json:"admin_addr"`
	Status      enum.ClusterNodeStatus `json:"status"`
}

type ClusterNodeInput struct {
	Source      string `json:"source"`
	ClusterAddr string `json:"cluster_addr"`
}
