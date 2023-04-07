package cluster_service

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	cluster_entry2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	cluster_model2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	cluster_store2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-store"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/eosc/common/bean"
	"gopkg.in/errgo.v2/errors"
	"time"
)

type clusterService struct {
	clusterStore           cluster_store2.IClusterStore
	clusterHistoryStore    cluster_store2.IClusterHistoryStore
	clusterRuntime         cluster_store2.IClusterRuntimeStore
	clusterNodeService     cluster.IClusterNodeService
	clusterVariableService variable.IClusterVariableService
	apintoClientService    cluster.IApintoClient
}

func newClusterService() cluster.IClusterService {
	s := &clusterService{}
	bean.Autowired(&s.clusterStore)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.clusterVariableService)
	bean.Autowired(&s.clusterHistoryStore)
	bean.Autowired(&s.clusterRuntime)
	bean.Autowired(&s.apintoClientService)

	return s
}

func (c *clusterService) GetAllCluster(ctx context.Context) ([]*cluster_model2.Cluster, error) {
	clusters, err := c.clusterStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.Cluster, 0, len(clusters))
	for _, clusterInfo := range clusters {
		item := &cluster_model2.Cluster{
			Cluster: clusterInfo,
		}
		list = append(list, item)
	}

	return list, nil
}

// CheckByNamespaceByName 检测集群是否存在 通过namespace和name检测
func (c *clusterService) CheckByNamespaceByName(ctx context.Context, namespaceId int, name string) (int, error) {
	clusterInfo, err := c.clusterStore.GetByNamespaceByName(ctx, namespaceId, name)
	if err != nil {
		return 0, common.ClusterNotExist
	}
	return clusterInfo.Id, nil
}

// GetByNamespaceByName 检测集群是否存在 通过namespace和name检测
func (c *clusterService) GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_model2.Cluster, error) {
	clusterInfo, err := c.clusterStore.GetByNamespaceByName(ctx, namespaceId, name)
	if err != nil {
		return nil, common.ClusterNotExist
	}

	//TODO 获取集群状态

	return &cluster_model2.Cluster{Cluster: clusterInfo}, nil
}

// GetByClusterId 根据集群ID获取集群信息
func (c *clusterService) GetByClusterId(ctx context.Context, clusterId int) (*cluster_model2.Cluster, error) {
	clusterInfo, err := c.clusterStore.Get(ctx, clusterId)
	if err != nil {
		return nil, errors.New("cluster does not exist")
	}

	value := &cluster_model2.Cluster{
		Cluster: clusterInfo,
	}
	return value, nil
}

func (c *clusterService) GetByNamespaceId(ctx context.Context, namespaceId int) ([]*cluster_model2.Cluster, error) {
	clusters, err := c.clusterStore.GetAllByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.Cluster, 0, len(clusters))
	for _, clusterInfo := range clusters {
		value := &cluster_model2.Cluster{
			Cluster: clusterInfo,
		}
		list = append(list, value)
	}

	return list, nil
}

func (c *clusterService) GetByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_model2.Cluster, error) {
	clusters, err := c.clusterStore.GetByNamespaceByNames(ctx, namespaceId, names)
	if err != nil {
		return nil, err
	}
	list := make([]*cluster_model2.Cluster, 0, len(clusters))
	for _, clusterInfo := range clusters {
		value := &cluster_model2.Cluster{
			Cluster: clusterInfo,
		}
		list = append(list, value)
	}

	return list, nil
}

// Insert 新增集群
func (c *clusterService) Insert(ctx context.Context, namespaceId, userId int, clusterInput *cluster_dto.ClusterInput) error {
	clusterId, _ := c.CheckByNamespaceByName(ctx, namespaceId, clusterInput.Name)
	if clusterId > 0 {
		return errors.New("cluster already exists")
	}

	bytes, err := common.Base64Decode(clusterInput.Source)
	if err != nil {
		return err
	}
	nodes := make([]*cluster_model2.ClusterNode, 0)
	if err = json.Unmarshal(bytes, &nodes); err != nil {
		return err
	}

	//节点重复对比
	if err = c.clusterNodeService.NodeRepeatContrast(ctx, namespaceId, clusterId, nodes); err != nil {
		return err
	}

	clusterInfo, err := c.clusterNodeService.GetClusterInfo(clusterInput.Addr)
	if err != nil {
		return err
	}

	t := time.Now()
	entryCluster := &cluster_entry2.Cluster{
		NamespaceId: namespaceId,
		Name:        clusterInput.Name,
		Desc:        clusterInput.Desc,
		Env:         clusterInput.Env,
		Addr:        clusterInput.Addr,
		UUID:        clusterInfo.Cluster,
		CreateTime:  t,
		UpdateTime:  t,
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: clusterInput.Name,
	})

	return c.clusterStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.clusterStore.Insert(txCtx, entryCluster); err != nil {
			return err
		}

		if err = c.clusterHistoryStore.HistoryAdd(txCtx, namespaceId, entryCluster.Id, *entryCluster, userId); err != nil {
			return err
		}

		entryClusterNodes := make([]*cluster_model2.ClusterNode, 0, len(nodes))
		nodesAdminAddr := make([]string, 0, len(nodes))

		for _, node := range nodes {
			entryClusterNodes = append(entryClusterNodes, &cluster_model2.ClusterNode{ClusterNode: &cluster_entry2.ClusterNode{
				NamespaceId: namespaceId,
				AdminAddr:   node.AdminAddr,
				ServiceAddr: node.ServiceAddr,
				Name:        node.Name,
				CreateTime:  t,
				ClusterId:   entryCluster.Id,
			}})
			nodesAdminAddr = append(nodesAdminAddr, node.AdminAddr)
		}
		for _, node := range entryClusterNodes {
			node.ClusterId = entryCluster.Id
		}

		return c.clusterNodeService.Insert(txCtx, entryClusterNodes)
	})
}

// QueryByNamespaceId 根据工作空间ID查询所有的集群
func (c *clusterService) QueryByNamespaceId(ctx context.Context, namespaceId int, clusterName string) (*cluster_model2.Cluster, error) {

	clusterId, err := c.CheckByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	return c.GetByClusterId(ctx, clusterId)

}

// QueryListByNamespaceId 根据工作空间ID查询所有的集群 提供给web端接口
func (c *clusterService) QueryListByNamespaceId(ctx context.Context, namespaceId int) ([]*cluster_model2.Cluster, error) {
	clusters, err := c.clusterStore.GetAllByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	list := make([]*cluster_model2.Cluster, 0, len(clusters))
	for _, clusterInfo := range clusters {
		status := 1
		clusterNodes, _, _ := c.clusterNodeService.QueryList(ctx, namespaceId, clusterInfo.Name)
		if len(clusterNodes) == 0 {
			status = 3 //异常
		} else {
			abnormalNum := 0
			normalNum := 0
			for _, node := range clusterNodes {
				if node.Status == 2 {
					normalNum++
				} else {
					abnormalNum++
				}
			}
			if normalNum == len(clusterNodes) { //正常
				status = 1
			} else if abnormalNum == len(clusterNodes) {
				status = 3
			} else {
				status = 2 //部分正常
			}
		}

		//兼容旧版本数据
		if clusterInfo.UUID == "" {
			go func() {
				info, _ := c.clusterNodeService.GetClusterInfo(clusterInfo.Addr)
				if info != nil {
					_ = c.UpdateAddr(ctx, 0, clusterInfo.Id, clusterInfo.Addr, clusterInfo.UUID)
				}
			}()
		}

		list = append(list, &cluster_model2.Cluster{
			Cluster: clusterInfo,
			Status:  status,
		})
	}

	return list, nil
}

// DeleteByNamespaceIdByName 删除集群
func (c *clusterService) DeleteByNamespaceIdByName(ctx context.Context, namespaceId, userId int, name string) error {
	clusterId, err := c.CheckByNamespaceByName(ctx, namespaceId, name)
	if err != nil {
		return err
	}
	clusterInfo, err := c.clusterStore.Get(ctx, clusterId)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: name,
	})

	return c.clusterStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = c.clusterStore.Delete(txCtx, clusterId); err != nil {
			return err
		}

		//删除runtime表中该集群的所有记录
		if err = c.clusterRuntime.DeleteByClusterID(txCtx, clusterId); err != nil {
			return err
		}

		//新增删除历史
		if err = c.clusterHistoryStore.HistoryDelete(txCtx, namespaceId, clusterId, *clusterInfo, userId); err != nil {
			return err
		}

		//删除集群下的环境变量
		if err = c.clusterVariableService.DeleteAll(txCtx, namespaceId, clusterId, userId); err != nil {
			return err
		}

		return nil
	})

}

// UpdateDesc 修改集群描述
func (c *clusterService) UpdateDesc(ctx context.Context, namespaceId, userId int, name, desc string) error {
	clusterId, err := c.CheckByNamespaceByName(ctx, namespaceId, name)
	if err != nil {
		return err
	}

	clusterInfo := &cluster_entry2.Cluster{
		Id:   clusterId,
		Desc: desc,
	}
	oldCluster, err := c.clusterStore.Get(ctx, clusterId)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: name,
	})

	return c.clusterStore.Transaction(ctx, func(txCtx context.Context) error {
		_, err = c.clusterStore.Update(txCtx, clusterInfo)
		if err != nil {
			return err
		}

		return c.clusterHistoryStore.HistoryEdit(txCtx, namespaceId, clusterId, oldCluster, clusterInfo, userId)
	})
}

func (c *clusterService) UpdateAddr(ctx context.Context, userId, clusterId int, addr, uuid string) error {
	clusterInfo := &cluster_entry2.Cluster{
		Id:   clusterId,
		Addr: addr,
		UUID: uuid,
	}
	oldCluster, err := c.clusterStore.Get(ctx, clusterId)
	if err != nil {
		return err
	}
	return c.clusterStore.Transaction(ctx, func(txCtx context.Context) error {
		_, err = c.clusterStore.Update(txCtx, clusterInfo)
		if err != nil {
			return err
		}

		return c.clusterHistoryStore.HistoryEdit(txCtx, oldCluster.NamespaceId, clusterInfo.Id, oldCluster, clusterInfo, userId)
	})

}
