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
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

type clusterNodeService struct {
	clusterNodeStore cluster_store.IClusterNodeStore
	clusterService   cluster.IClusterService
	apintoClient     cluster.IApintoClient
}

func newClusterNodeService() cluster.IClusterNodeService {
	s := &clusterNodeService{}
	bean.Autowired(&s.clusterNodeStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.apintoClient)

	return s
}

func (c *clusterNodeService) Insert(ctx context.Context, nodes []*cluster_model2.ClusterNode) error {
	entryNodes := make([]*cluster_entry.ClusterNode, 0, len(nodes))
	for _, node := range nodes {
		entryNodes = append(entryNodes, node.ClusterNode)
	}
	return c.clusterNodeStore.Insert(ctx, entryNodes...)
}

// QueryList 查询集群下的节点列表
func (c *clusterNodeService) QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model2.ClusterNode, bool, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, false, err
	}

	list := make([]*cluster_model2.ClusterNode, 0)

	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		nodes, err := c.clusterNodeStore.GetAllByClusterIds(ctx, cluster.Id)
		if err != nil {
			return err
		}
		for _, node := range nodes {
			list = append(list, &cluster_model2.ClusterNode{
				ClusterNode: node,
				AdminAddrs:  strings.Split(node.AdminAddr, ","),
				Status:      1,
			})
		}
		return nil
	})

	clusterNodes := make([]*cluster_model2.ClusterNode, 0)
	group.Go(func() error {
		clusterNodes, err = c.GetNodesByUrl(cluster.Addr)
		if err != nil {
			log.Errorf("clusterNodeService-QueryList addr=%s err=%s", cluster.Addr, err.Error())
			return err
		}
		return nil
	})

	if err = group.Wait(); err != nil {
		return nil, false, err
	}

	if len(clusterNodes) > 0 {
		for _, node := range list {
			node.Status = 2
		}
	}

	addList, updateList, delList := common.DiffContrast(list, clusterNodes)

	isUpdate := len(addList) > 0 || len(updateList) > 0 || len(delList) > 0

	return list, isUpdate, nil
}

func (c *clusterNodeService) QueryByClusterIds(ctx context.Context, clusterIds ...int) ([]*cluster_model2.ClusterNode, error) {
	nodes, err := c.clusterNodeStore.GetAllByClusterIds(ctx, clusterIds...)
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.ClusterNode, 0, len(nodes))

	for _, node := range nodes {
		list = append(list, &cluster_model2.ClusterNode{
			ClusterNode: node,
			AdminAddrs:  strings.Split(node.AdminAddr, ","),
		})
	}
	return list, nil
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
	nodes := make([]*cluster_model2.ClusterNode, 0)
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

	err = c.clusterNodeStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		entryClusterNodes := make([]*cluster_entry.ClusterNode, 0, len(nodes))
		nodesAdminAddr := make([]string, 0, len(nodes))

		for _, node := range nodes {
			entryClusterNodes = append(entryClusterNodes, &cluster_entry.ClusterNode{
				ClusterId:   clusterId,
				NamespaceId: namespaceId,
				AdminAddr:   node.AdminAddr,
				ServiceAddr: node.ServiceAddr,
				Name:        node.Name,
				CreateTime:  t,
			})
			nodesAdminAddr = append(nodesAdminAddr, node.AdminAddr)
		}
		if err = c.clusterNodeStore.UpdateNodes(txCtx, clusterId, entryClusterNodes); err != nil {
			return err
		}

		return c.clusterService.UpdateAddr(txCtx, userId, clusterId, clusterAddr, clusterInfo.Cluster)
	})
	if err != nil {
		return err
	}
	c.apintoClient.SetClient(namespaceId, clusterId)
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
			AdminAddr:   node.AdminAddr,
			ServiceAddr: node.ServiceAddr,
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
	c.apintoClient.SetClient(namespaceId, clusterInfo.Id)
	return nil
}

// NodeRepeatContrast  节点重复对比
func (c *clusterNodeService) NodeRepeatContrast(ctx context.Context, namespaceId, clusterId int, newList []*cluster_model2.ClusterNode) error {
	clusters, err := c.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	clustersMap := common.SliceToMap(clusters, func(t *cluster_model2.Cluster) int {
		return t.Id
	})
	//工作空间下任何一个节点名称和这次添加的有重复,不可保存
	clusterIds := make([]int, 0)
	for _, clusterInfo := range clusters {
		if clusterId == clusterInfo.Id { //过滤本身的
			continue
		}
		clusterIds = append(clusterIds, clusterInfo.Id)
	}

	clusterNodes, err := c.QueryByClusterIds(ctx, clusterIds...)
	if err != nil {
		return err
	}

	//对比clusterNods和nodes是否有重复的name

	mapNode := common.SliceToMap(newList, func(t *cluster_model2.ClusterNode) string {
		return t.Name
	})

	for _, node := range clusterNodes {
		if _, ok := mapNode[node.Name]; ok {
			if clusterInfo, clusterOk := clustersMap[node.ClusterId]; clusterOk {
				return errors.New(fmt.Sprintf("%s集群已有这个节点信息", clusterInfo.Name))
			}
			return errors.New(fmt.Sprintf("%s集群已有这个节点信息", node.Name))
		}
	}

	return nil
}

func (c *clusterNodeService) GetNodesByUrl(addr string) ([]*cluster_model2.ClusterNode, error) {

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

	list := make([]*cluster_model2.ClusterNode, 0)
	for _, node := range clusterInfo.Nodes {
		adminAddr := strings.Join(node.Admin, ",")
		serverAddr := strings.Join(node.Server, ",")
		list = append(list, &cluster_model2.ClusterNode{
			ClusterNode: &cluster_entry.ClusterNode{
				Name:        node.Name,
				AdminAddr:   adminAddr,
				ServiceAddr: serverAddr,
			},
			Status: 2,
		})
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
