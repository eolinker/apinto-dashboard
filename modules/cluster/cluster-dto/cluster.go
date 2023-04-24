package cluster_dto

import "github.com/eolinker/apinto-dashboard/enum"

type ClusterInput struct {
	Id     int    `json:"id"`
	Name   string `json:"name,omitempty"`
	Env    string `json:"env,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Source string `json:"source,omitempty"`
}

// ClusterOut 集群信息
type ClusterOut struct {
	Id         int                `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Env        string             `json:"env,omitempty"`
	Status     enum.ClusterStatus `json:"status,omitempty"`
	UUID       string             `json:"uuid,omitempty"`
	Desc       string             `json:"desc,omitempty"`
	CreateTime string             `json:"create_time,omitempty"`
	UpdateTime string             `json:"update_time,omitempty"`
}

type EnvCluster struct {
	Clusters []*ClusterOut `json:"clusters"`
	Name     string        `json:"name"`
}
