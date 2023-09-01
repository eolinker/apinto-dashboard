package cluster_service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/eosc/common/bean"
)

type apintoClientService struct {
	clusterNodeService cluster.IClusterNodeService
	clusterService     cluster.IClusterService
}

func newApintoClientService() cluster.IApintoClient {
	s := &apintoClientService{}
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.clusterNodeService)
	//bean.Autowired(&s.resetOnline)
	return s
}

func (c *apintoClientService) GetClient(ctx context.Context, clusterId int) (v1.IClient, error) {
	adminAddr, err := c.clusterNodeService.QueryAdminAddrByClusterId(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	client, err := v1.NewClient(adminAddr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InitClustersGlobalPlugin 初始化所有集群的全局插件配置
//func (c *apintoClientService) InitClustersGlobalPlugin(ctx context.Context) error {
//	//获取所有的集群id
//	clusters, err := c.clusterService.GetAllCluster(ctx)
//	if err != nil {
//		return err
//	}
//
//	pluginList := plugin.GetGlobalPluginConf()
//
//	//判断摘要不同，则进行更新?
//
//	for _, cluster := range clusters {
//		//连接得上的正常节点才进行更新
//		client, err := c.getClient(ctx, cluster.Id)
//		if err != nil {
//			log.Infof("cluster %s abnormal. Can't init globalPlugin.", cluster.Name)
//			continue
//		}
//		_ = client.ForGlobalPlugin().Set(pluginList)
//	}
//
//	return nil
//}

//func (c *apintoClientService) InitGlobalPlugin(clusterAddr string, nodesAdminAddr []string) error {
//	//newAddrSlice := strings.SplitN(strings.ReplaceAll(clusterAddr, "http://", ""), ".", 3)
//	newAdmin := make([]string, 0)
//	//newAdmin = append(newAdmin, clusterAddr)
//	for _, adminAddr := range nodesAdminAddr {
//
//		for _, nodeAddr := range strings.Split(adminAddr, ",") {
//			//newNodeAddrSlice := strings.SplitN(strings.ReplaceAll(nodeAddr, "http://", ""), ".", 3)
//			//if len(newNodeAddrSlice) >= 2 {
//			//	if newAddrSlice[0] == newNodeAddrSlice[0] && newAddrSlice[1] == newNodeAddrSlice[1] {
//			newAdmin = append(newAdmin, nodeAddr)
//			//}
//			//}
//		}
//
//	}
//	client, err := v1.NewClient(newAdmin)
//	if err != nil {
//		return fmt.Errorf("cluster init global plugin fail. err:%s. ", err)
//	}
//
//	pluginList := plugin.GetGlobalPluginConf()
//
//	err = client.ForGlobalPlugin().Set(pluginList)
//	if err != nil {
//		return fmt.Errorf("cluster init global plugin fail. err:%s. ", err)
//	}
//
//	return nil
//}
