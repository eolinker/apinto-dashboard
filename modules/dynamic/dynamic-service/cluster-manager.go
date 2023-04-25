package dynamic_service

import (
	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/eosc"
)

var (
	clusterManager = NewClusterManager()
)

type IClusterManager interface {
	GetClient(cluster string, addr string) v2.IClient
	ClusterKeys() []string
}

func NewClusterManager() IClusterManager {
	return &ClusterManager{
		clusterClients: eosc.BuildUntyped[string, v2.IClient](),
	}
}

type ClusterManager struct {
	clusterClients eosc.Untyped[string, v2.IClient]
}

func (c *ClusterManager) ClusterKeys() []string {
	return c.clusterClients.Keys()
}

func (c *ClusterManager) GetClient(cluster string, addr string) v2.IClient {
	client, has := c.clusterClients.Get(cluster)
	if has {
		if client.Addr() == addr {
			return client
		}
	}
	newClient := v2.NewClient(addr)
	c.clusterClients.Set(cluster, newClient)
	return newClient
}

func (c *ClusterManager) Remove(clusters []string) {
	c.clusterClients.Dels(clusters...)
}

func GetClusterClient(cluster string, addr string) v2.IClient {
	return clusterManager.GetClient(cluster, addr)
}
