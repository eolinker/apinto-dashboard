package cluster_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	cluster_model2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type clusterNodeService struct {
	clusterNodeStore cluster_store.IClusterNodeStore
	clusterService   cluster.IClusterService
	nodeCache        INodeCache
	apintoClient     cluster.IApintoClient
}

func (c *clusterNodeService) List(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model2.Node, error) {
	cm, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}
	return c.QueryByClusterId(ctx, cm.Id)

}

func (c *clusterNodeService) Delete(ctx context.Context, namespaceId int, clusterId int) error {

	_, err := c.clusterNodeStore.DeleteWhere(ctx, map[string]interface{}{
		"namespace": namespaceId,
		"cluster":   clusterId,
	})
	if err != nil {
		return err
	}
	return nil
}

func newClusterNodeService() cluster.IClusterNodeService {
	s := &clusterNodeService{}
	bean.Autowired(&s.clusterNodeStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.nodeCache)

	return s
}

func (c *clusterNodeService) Insert(ctx context.Context, nodes []*cluster_model2.Node) error {
	entryNodes := make([]*cluster_entry.ClusterNode, 0, len(nodes))
	for _, node := range nodes {
		entryNodes = append(entryNodes, node.ToEntity())
	}
	return c.clusterNodeStore.Insert(ctx, entryNodes...)
}

// QueryList 查询集群下的节点列表
func (c *clusterNodeService) QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model2.ClusterNode, bool, error) {
	clusterModel, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, false, err
	}

	//控制台集群下存储的节点列表
	list := make([]*cluster_model2.Node, 0)

	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		nodes, err := c.QueryByClusterId(ctx, clusterModel.Id)
		if err != nil {
			return err
		}
		list = append(list, nodes...)
		//for _, node := range nodes {
		//	list = append(list, &cluster_model2.Node{
		//		Node:   *node,
		//		Status: 1,
		//	})
		//}
		return nil
	})
	//调用admin接口返回集群下的节点列表
	clusterNodes := make([]*cluster_model2.Node, 0)
	group.Go(func() error {
		ns, err := c.GetNodesByUrl(clusterModel.Addr)
		if err != nil {
			log.Errorf("clusterNodeService-QueryList addr=%s err=%s", clusterModel.Addr, err.Error())
			return err
		}
		clusterNodes = append(clusterNodes, ns...)
		//for _, n := range ns {
		//	clusterNodes = append(clusterNodes, &cluster_model2.ClusterNode{
		//		Node:   *n,
		//		Status: 1,
		//	})
		//}
		return nil
	})

	if err = group.Wait(); err != nil {
		return nil, false, err
	}

	rs := make([]*cluster_model2.ClusterNode, 0, len(list))
	addList, updateList, delList := common.DiffContrast(list, clusterNodes)
	isUpdate := len(addList) > 0 || len(updateList) > 0 || len(delList) > 0

	if len(clusterNodes) > 0 {
		//判断ClusterNodes 集合里 是否存在控制台节点，
		clusterNodesSet := common.SliceToSet(clusterNodes, func(t *cluster_model2.Node) string {
			return t.Name
		})
		for _, node := range list {
			ni := &cluster_model2.ClusterNode{Node: *node, Status: 1}
			if _, has := clusterNodesSet[node.Name]; has {
				ni.Status = 2
			}
			rs = append(rs, ni)
		}
	}

	return rs, isUpdate, nil
}

func (c *clusterNodeService) QueryByClusterId(ctx context.Context, id int) ([]*cluster_model2.Node, error) {
	nodes, err := c.nodeCache.Get(ctx, id)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	//若redis存在值
	if err == nil {
		return *nodes, nil
	}

	//若redis无缓存
	nodeEntries, err := c.clusterNodeStore.GetAllByClusterIds(ctx, id)
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.Node, 0, len(nodeEntries))
	for _, node := range nodeEntries {
		list = append(list, cluster_model2.ReadClusterNode(node))
	}
	//缓存
	_ = c.nodeCache.Set(ctx, id, &list)
	return list, nil

}

func (c *clusterNodeService) QueryAllCluster(ctx context.Context) ([]*cluster_model2.Node, error) {
	nodes, err := c.clusterNodeStore.List(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.Node, 0, len(nodes))

	for _, node := range nodes {
		list = append(list, cluster_model2.ReadClusterNode(node))
	}
	return list, nil
}

func (c *clusterNodeService) QueryAdminAddrByClusterId(ctx context.Context, id int) ([]string, error) {
	nodes, err := c.QueryByClusterId(ctx, id)
	if err != nil {
		return nil, err
	}
	admins := make([]string, 0)
	for _, node := range nodes {
		for _, nodeAddr := range node.AdminAddrs {
			admins = append(admins, nodeAddr)
		}
	}
	return admins, nil
}

// Reset 重置节点信息
func (c *clusterNodeService) Reset(ctx context.Context, namespaceId, userId int, clusterName, clusterAddr, source string) error {
	clusterId, err := c.clusterService.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	bytes, err := common.Base64Decode(source)
	if err != nil {
		return err
	}
	nodes := make([]*cluster_model2.Node, 0)
	if err = json.Unmarshal(bytes, &nodes); err != nil {
		return err
	}

	//节点重复对比
	if err = c.NodeRepeatContrast(ctx, namespaceId, clusterId, nodes); err != nil {
		return err
	}

	clusterInfo, err := c.GetClusterInfo(clusterAddr)
	if err != nil {
		return err
	}

	entryClusterNodes := make([]*cluster_entry.ClusterNode, 0, len(nodes))
	err = c.clusterNodeStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()

		nodesAdminAddr := make([]string, 0, len(nodes))

		for _, node := range nodes {
			ent := node.ToEntity()
			ent.ClusterId = clusterId
			ent.CreateTime = t
			ent.NamespaceId = namespaceId

			entryClusterNodes = append(entryClusterNodes, ent)
			nodesAdminAddr = append(nodesAdminAddr, ent.AdminAddr)
		}
		if err = c.clusterNodeStore.UpdateNodes(txCtx, clusterId, entryClusterNodes); err != nil {
			return err
		}

		return c.clusterService.UpdateAddr(txCtx, userId, clusterId, clusterAddr, clusterInfo.Cluster)
	})
	if err != nil {
		return err
	}
	list := make([]*cluster_model2.Node, 0, len(entryClusterNodes))
	for _, node := range entryClusterNodes {
		list = append(list, cluster_model2.ReadClusterNode(node))
	}
	//缓存
	_ = c.nodeCache.Set(ctx, clusterId, &list)

	return nil
}

// Update 更新节点信息
func (c *clusterNodeService) Update(ctx context.Context, namespaceId int, clusterName string) error {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	t := time.Now()
	nodes, err := c.GetNodesByUrl(clusterInfo.Addr)
	if err != nil {
		return err
	}
	//节点重复对比
	if err = c.NodeRepeatContrast(ctx, namespaceId, clusterInfo.Id, nodes); err != nil {
		return err
	}

	//由于节点名称是可能会变的，所以重新更新节点是把原来的全部删除，然后重新添加新的节点信息。
	newClusterNodes := make([]*cluster_entry.ClusterNode, 0, len(nodes))
	for _, node := range nodes {
		newClusterNodes = append(newClusterNodes, &cluster_entry.ClusterNode{
			ClusterId:   clusterInfo.Id,
			NamespaceId: namespaceId,
			AdminAddr:   strings.Join(node.AdminAddrs, ","),
			ServiceAddr: strings.Join(node.ServiceAddr, ","),
			Name:        node.Name,
			CreateTime:  t,
		})
	}

	err = c.clusterNodeStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.clusterNodeStore.UpdateNodes(txCtx, clusterInfo.Id, newClusterNodes); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	list := make([]*cluster_model2.Node, 0, len(newClusterNodes))
	for _, node := range nodes {
		node.NamespaceId = namespaceId
		node.ClusterId = clusterInfo.Id
		node.CreateTime = t
		list = append(list, node)
	}
	//缓存
	_ = c.nodeCache.Set(ctx, clusterInfo.Id, &list)

	return nil
}

// NodeRepeatContrast  节点重复对比
func (c *clusterNodeService) NodeRepeatContrast(ctx context.Context, namespaceId, clusterId int, newList []*cluster_model2.Node) error {
	clusters, err := c.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	clustersMap := common.SliceToMap(clusters, func(t *cluster_model2.Cluster) int {
		return t.Id
	})

	clusterNodes, err := c.QueryAllCluster(ctx)
	if err != nil {
		return err
	}
	//工作空间下任何一个节点名称和这次添加的有重复,不可保存
	filteredNodes := make([]*cluster_model2.Node, 0, len(clusterNodes))
	for _, node := range clusterNodes {
		if clusterId == node.ClusterId { //过滤本身的
			continue
		}
		filteredNodes = append(filteredNodes, node)
	}

	//对比clusterNods和nodes是否有重复的name

	mapNode := common.SliceToMap(newList, func(t *cluster_model2.Node) string {
		return t.Name
	})

	for _, node := range filteredNodes {
		if _, ok := mapNode[node.Name]; ok {
			if clusterInfo, clusterOk := clustersMap[node.ClusterId]; clusterOk {
				return errors.New(fmt.Sprintf("%s集群已有这个节点信息", clusterInfo.Name))
			}
			return errors.New(fmt.Sprintf("%s集群已有这个节点信息", node.Name))
		}
	}

	return nil
}

func (c *clusterNodeService) GetNodesByUrl(addr string) ([]*cluster_model2.Node, error) {

	client, err := v1.NewClient([]string{addr})
	if err != nil {
		return nil, err
	}

	newAddrSlice := strings.SplitN(strings.ReplaceAll(addr, "http://", ""), ".", 3)

	if len(newAddrSlice) > 0 {
		if newAddrSlice[0] == "127" {
			return nil, errors.New("不能使用本地地址")
		}
	}

	clusterInfo, err := client.ClusterInfo()
	if err != nil {
		return nil, err
	}

	list := make([]*cluster_model2.Node, 0)
	for _, node := range clusterInfo.Nodes {
		list = append(list,

			&cluster_model2.Node{
				Name:        node.Name,
				AdminAddrs:  node.Admin,
				ServiceAddr: node.Server,
			},
		)
	}

	return list, nil
}

func (c *clusterNodeService) GetClusterInfo(addr string) (*v1.ClusterInfo, error) {
	client, err := v1.NewClient([]string{addr})
	if err != nil {
		return nil, err
	}

	clusterInfo, err := client.ClusterInfo()
	if err != nil {
		return nil, err
	}

	return clusterInfo, nil
}
