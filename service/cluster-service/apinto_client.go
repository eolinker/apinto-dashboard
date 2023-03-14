package cluster_service

import (
	"context"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/client/v1/initialize/plugin"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"strings"
	"sync"
)

type IApintoClient interface {
	GetClient(ctx context.Context, clusterId int) (v1.IClient, error)
	SetClient(namespace, clusterId int)
	InitClustersGlobalPlugin(ctx context.Context) error
	InitGlobalPlugin(clusterAddr string, nodesAdminAddr []string) error
}

type apintoClientService struct {
	clusterNodeService IClusterNodeService
	clusterService     IClusterService
	resetOnline        IResetOnlineService
	lock               *sync.Mutex
	clientMap          map[int]v1.IClient
}

func (c *apintoClientService) SetClient(namespaceId, clusterId int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.clientMap, clusterId)
	//重新上线
	go c.resetOnline.ResetOnline(context.Background(), namespaceId, clusterId)
}

func newApintoClientService() IApintoClient {
	s := &apintoClientService{
		lock:      new(sync.Mutex),
		clientMap: make(map[int]v1.IClient),
	}
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.resetOnline)
	return s
}

func (c *apintoClientService) GetClient(ctx context.Context, clusterId int) (v1.IClient, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if v, ok := c.clientMap[clusterId]; ok {
		return v, nil
	}
	client, err := c.getClient(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	c.clientMap[clusterId] = client
	return client, nil
}

func (c *apintoClientService) getClient(ctx context.Context, clusterId int) (v1.IClient, error) {
	nodes, err := c.clusterNodeService.QueryByClusterIds(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	//cluster, err := c.clusterService.GetByClusterId(ctx, clusterId)
	//if err != nil {
	//	return nil, err
	//}

	//newAddrSlice := strings.SplitN(strings.ReplaceAll(cluster.Addr, "http://", ""), ".", 3)

	newAdmin := make([]string, 0)
	//newAdmin = append(newAdmin, cluster.Addr)
	for _, node := range nodes {

		for _, nodeAddr := range strings.Split(node.AdminAddr, ",") {
			//newNodeAddrSlice := strings.SplitN(strings.ReplaceAll(nodeAddr, "http://", ""), ".", 3)
			//if len(newNodeAddrSlice) >= 2 {
			//	if newAddrSlice[0] == newNodeAddrSlice[0] && newAddrSlice[1] == newNodeAddrSlice[1] {
			newAdmin = append(newAdmin, nodeAddr)
			//}
			//}
		}

	}

	client, err := v1.NewClient(newAdmin)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InitClustersGlobalPlugin 初始化所有集群的全局插件配置
func (c *apintoClientService) InitClustersGlobalPlugin(ctx context.Context) error {
	//获取所有的集群id
	clusters, err := c.clusterService.GetAllCluster(ctx)
	if err != nil {
		return err
	}

	pluginList := plugin.GetGlobalPluginConf()

	//判断摘要不同，则进行更新?

	for _, cluster := range clusters {
		//连接得上的正常节点才进行更新
		client, err := c.getClient(ctx, cluster.Id)
		if err != nil {
			log.Infof("cluster %s abnormal. Can't init globalPlugin.", cluster.Name)
			continue
		}
		_ = client.ForGlobalPlugin().Set(pluginList)
	}

	return nil
}

func (c *apintoClientService) InitGlobalPlugin(clusterAddr string, nodesAdminAddr []string) error {
	//newAddrSlice := strings.SplitN(strings.ReplaceAll(clusterAddr, "http://", ""), ".", 3)
	newAdmin := make([]string, 0)
	//newAdmin = append(newAdmin, clusterAddr)
	for _, adminAddr := range nodesAdminAddr {

		for _, nodeAddr := range strings.Split(adminAddr, ",") {
			//newNodeAddrSlice := strings.SplitN(strings.ReplaceAll(nodeAddr, "http://", ""), ".", 3)
			//if len(newNodeAddrSlice) >= 2 {
			//	if newAddrSlice[0] == newNodeAddrSlice[0] && newAddrSlice[1] == newNodeAddrSlice[1] {
			newAdmin = append(newAdmin, nodeAddr)
			//}
			//}
		}

	}
	client, err := v1.NewClient(newAdmin)
	if err != nil {
		return fmt.Errorf("cluster init global plugin fail. err:%s. ", err)
	}

	pluginList := plugin.GetGlobalPluginConf()

	err = client.ForGlobalPlugin().Set(pluginList)
	if err != nil {
		return fmt.Errorf("cluster init global plugin fail. err:%s. ", err)
	}

	return nil
}
