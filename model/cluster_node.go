package model

import (
	"github.com/eolinker/apinto-dashboard/entry/cluster-entry"
)

type ClusterNode struct {
	*cluster_entry.ClusterNode
	AdminAddrs []string
	Status     int //状态（1未运行,2运行中）
}

func (c *ClusterNode) GetKey() string {
	return c.Name
}

func (c *ClusterNode) Values() []string {
	return []string{c.AdminAddr}
}
