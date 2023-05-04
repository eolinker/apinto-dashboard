package cluster_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/eolinker/apinto-dashboard/controller/users"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type clusterController struct {
	clusterService     cluster.IClusterService
	clusterNodeService cluster.IClusterNodeService
}

func newClusterController() *clusterController {
	c := &clusterController{}
	bean.Autowired(&c.clusterService)
	bean.Autowired(&c.clusterNodeService)
	return c
}

// clusters 获取集群列表
func (c *clusterController) clusters(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	clusters, err := c.clusterService.QueryListByNamespaceId(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	list := make([]*cluster_dto.ClusterOut, 0, len(clusters))
	for _, clusterInfo := range clusters {
		list = append(list, &cluster_dto.ClusterOut{
			Name:       clusterInfo.Name,
			Title:      clusterInfo.Title,
			Env:        clusterInfo.Env,
			Status:     enum.ClusterStatus(clusterInfo.Status),
			Desc:       clusterInfo.Desc,
			CreateTime: common.TimeToStr(clusterInfo.CreateTime),
			UpdateTime: common.TimeToStr(clusterInfo.UpdateTime),
		})
	}

	m := common.Map[string, interface{}]{}

	m["clusters"] = list

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))

}

// simpleClusters 获取简易集群列表
func (c *clusterController) simpleClusters(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	clusters, err := c.clusterService.SimpleCluster(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"clusters": clusters,
	}))

}

func (c *clusterController) clusterEnum(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	list, err := c.clusterService.GetByNamespaceId(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	envCluster := make([]*cluster_dto.EnvCluster, 0)
	clusterMap := common.Map[string, []*cluster_model.Cluster]{}

	for _, clusterInfo := range list {
		clusterMap[clusterInfo.Env] = append(clusterMap[clusterInfo.Env], clusterInfo)
	}

	for env, clusters := range clusterMap {
		clusterOuts := make([]*cluster_dto.ClusterOut, 0)

		for _, clusterInfo := range clusters {
			clusterOuts = append(clusterOuts, &cluster_dto.ClusterOut{
				Name: clusterInfo.Name,
				UUID: clusterInfo.UUID,
			})
		}

		envCluster = append(envCluster, &cluster_dto.EnvCluster{
			Clusters: clusterOuts,
			Name:     env,
		})
	}

	sort.Slice(envCluster, func(i, j int) bool {
		return envCluster[i].Name > envCluster[j].Name
	})
	m := common.Map[string, interface{}]{}
	m["envs"] = envCluster
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// get 获取单个集群信息
func (c *clusterController) cluster(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "Get cluster info fail. err: cluster_name can't be nil.")
		return
	}

	clusterInfo, err := c.clusterService.QueryByNamespaceId(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	m := common.Map[string, interface{}]{}
	m["cluster"] = &cluster_dto.ClusterOut{
		Name:       clusterInfo.Name,
		Title:      clusterInfo.Title,
		Env:        clusterInfo.Env,
		Status:     enum.ClusterStatus(clusterInfo.Status),
		Desc:       clusterInfo.Desc,
		CreateTime: common.TimeToStr(clusterInfo.CreateTime),
		UpdateTime: common.TimeToStr(clusterInfo.UpdateTime),
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))

}

// create 新建集群
func (c *clusterController) create(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	input := &cluster_dto.ClusterInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//校验是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Name == "" || input.Env == "" || input.Addr == "" || input.Source == "" {
		fmt.Println(*input)
		controller.ErrorJson(ginCtx, http.StatusOK, "parameter error")
		return
	}
	userId := users.GetUserId(ginCtx)
	if err := c.clusterService.Insert(ginCtx, namespaceId, userId, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// del 删除集群
func (c *clusterController) del(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "Delete cluster fail. err: cluster_name can't be nil.")
		return
	}
	userId := users.GetUserId(ginCtx)
	if err := c.clusterService.DeleteByNamespaceIdByName(ginCtx, namespaceId, userId, clusterName); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// update 更新集群信息
func (c *clusterController) update(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	clusterInput := &cluster_dto.ClusterInput{}
	err := ginCtx.BindJSON(clusterInput)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	userId := users.GetUserId(ginCtx)
	if err = c.clusterService.Update(ginCtx, namespaceId, userId, clusterName, clusterInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

//// putDesc 修改集群描述
//func (c *clusterController) putDesc(ginCtx *gin.Context) {
//
//	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
//	clusterName := ginCtx.Param("cluster_name")
//
//	clusterInput := &cluster_dto.ClusterInput{}
//	err := ginCtx.BindJSON(clusterInput)
//	if err != nil {
//		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
//		return
//	}
//	userId := users.GetUserId(ginCtx)
//	if err = c.clusterService.UpdateDesc(ginCtx, namespaceId, userId, clusterName, clusterInput.Desc); err != nil {
//		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
//		return
//	}
//	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
//
//}

// test 集群测试按钮
func (c *clusterController) test(context *gin.Context) {

	clusterAddr := context.Query("cluster_addr")
	nodes, err := c.clusterNodeService.GetNodesByUrl(clusterAddr)

	if err != nil {
		context.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	list := make([]*cluster_dto.ClusterNode, 0, len(nodes))

	isUpdate := false
	for _, node := range nodes {
		status := enum.ClusterNodeStatus(node.Status)
		if status == enum.ClusterNodeStatusRunning {
			isUpdate = true
		}
		list = append(list, &cluster_dto.ClusterNode{
			Name:        node.Name,
			ServiceAddr: node.ServiceAddr,
			AdminAddr:   node.AdminAddr,
			Status:      status,
		})

	}
	bytes, _ := json.Marshal(nodes)
	source := common.Base64Encode(bytes)

	m := common.Map[string, interface{}]{}
	m["nodes"] = list
	m["source"] = source
	m["is_update"] = isUpdate

	context.JSON(http.StatusOK, controller.NewSuccessResult(m))

}
