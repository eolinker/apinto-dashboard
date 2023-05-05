package v2

import (
	"fmt"

	"github.com/eolinker/eosc"
)

var (
	clusterManager = NewClusterManager()
)

type IClusterManager interface {
	GetClient(cluster string, addr string) (IClient, error)
	ClusterKeys() []string
}

func NewClusterManager() IClusterManager {
	return &ClusterManager{
		clusterClients: eosc.BuildUntyped[string, IClient](),
	}
}

type ClusterManager struct {
	clusterClients eosc.Untyped[string, IClient]
}

func (c *ClusterManager) ClusterKeys() []string {
	return c.clusterClients.Keys()
}

func (c *ClusterManager) GetClient(cluster string, addr string) (IClient, error) {
	client, has := c.clusterClients.Get(cluster)
	if has {
		if client.Addr() == addr {
			return client, nil
		}
	}
	newClient := NewClient(addr)
	if err := newClient.Ping(); err != nil {
		return nil, err
	}
	c.clusterClients.Set(cluster, newClient)
	return newClient, nil
}

func (c *ClusterManager) Remove(clusters []string) {
	c.clusterClients.Dels(clusters...)
}

func GetClusterClient(cluster string, addr string) (IClient, error) {
	return clusterManager.GetClient(cluster, addr)
}

//func MergeClusterStatus(cluster string, profession string, addr string, currentVersions map[string]string) map[string]map[string]string {
//	versionMap := make(map[string]string)
//	clusterNames := make([]string, 0, len(clusters))
//	for name, addr := range clusters {
//		clusterNames = append(clusterNames, name)
//		client, err := clusterManager.GetClient(name, addr)
//		if err != nil {
//			log.Errorf("get client(%s) error: %w.", addr, err)
//			continue
//		}
//		workers, err := client.List(profession)
//		if err != nil {
//			log.Errorf("get worker(%s) list error: %w.", profession, err)
//			continue
//		}
//		for _, w := range workers {
//			versionMap[toVersionKey(w.BasicInfo.Name, name)] = w.BasicInfo.Version
//		}
//	}
//	result := make(map[string]map[string]string)
//	for name, version := range currentVersions {
//		clusterStatus := map[string]string{}
//		for _, clusterName := range clusterNames {
//			if v, ok := versionMap[toVersionKey(name, clusterName)]; ok {
//				if v != version {
//					clusterStatus[name] = statusPre
//				} else {
//					clusterStatus[name] = statusOnline
//				}
//				continue
//			}
//			clusterStatus[name] = statusOffline
//		}
//		result[name] = clusterStatus
//	}
//	return result
//}

//func GetStatus(cluster string, addr string, profession, name string, version string) (string, error) {
//	client, err := clusterManager.GetClient(cluster, addr)
//	if err != nil {
//		return statusOffline, fmt.Errorf("get client(%s) error: %w", addr, err)
//	}
//	info, err := client.Info(profession, name)
//	if err != nil {
//		return statusOffline, fmt.Errorf("get worker(%s@%s) info error: %w", name, profession, err)
//	}
//	status := statusPre
//	if info.BasicInfo.Version == version {
//		status = statusOnline
//	}
//	return status, nil
//}

func Online(cluster string, addr string, profession, name string, body *WorkerInfo[BasicInfo]) error {
	client, err := clusterManager.GetClient(cluster, addr)
	if err != nil {
		return fmt.Errorf("get client(%s) error: %w", addr, err)
	}
	return client.Set(profession, name, body)
}

func Offline(cluster string, addr string, profession, name string) error {
	client, err := clusterManager.GetClient(cluster, addr)
	if err != nil {
		return fmt.Errorf("get client(%s) error: %w", addr, err)
	}
	return client.Delete(profession, name)
}
