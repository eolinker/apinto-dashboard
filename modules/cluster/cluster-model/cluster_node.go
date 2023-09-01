package cluster_model

import (
	cluster_entry "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"strings"
	"time"
)

type Node struct {
	Id          int
	NamespaceId int
	ClusterId   int
	Name        string
	ServiceAddr []string
	CreateTime  time.Time
	AdminAddrs  []string
}

func ReadClusterNode(n *cluster_entry.ClusterNode) *Node {
	return &Node{
		Id:          n.Id,
		NamespaceId: n.NamespaceId,
		ClusterId:   n.ClusterId,
		Name:        n.Name,
		ServiceAddr: strings.Split(n.ServiceAddr, ","),
		CreateTime:  n.CreateTime,
		AdminAddrs:  strings.Split(n.AdminAddr, ","),
	}
}
func (n *Node) ToEntity() *cluster_entry.ClusterNode {
	return &cluster_entry.ClusterNode{
		Id:          n.Id,
		NamespaceId: n.NamespaceId,
		ClusterId:   n.ClusterId,
		Name:        n.Name,
		AdminAddr:   strings.Join(n.AdminAddrs, ""),
		ServiceAddr: strings.Join(n.ServiceAddr, ","),
		CreateTime:  time.Time{},
	}
}

type ClusterNode struct {
	Node
	Status int //状态（1未运行,2运行中）
}

func (n *Node) GetKey() string {
	return n.Name
}

func (n *Node) Values() []string {
	return n.AdminAddrs
}
