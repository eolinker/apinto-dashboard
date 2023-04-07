package cluster_controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type clusterNodeController struct {
	clusterNodeService cluster.IClusterNodeService
}

func newClusterNodeController() *clusterNodeController {
	c := &clusterNodeController{}
	bean.Autowired(&c.clusterNodeService)
	return c
}

func RegisterClusterNodeRouter(router gin.IRoutes) {
	c := newClusterNodeController()

	router.GET("/cluster/:cluster_name/nodes", c.nodes)
	router.POST("/cluster/:cluster_name/node/reset", c.reset)
	router.PUT("/cluster/:cluster_name/node", c.put)
}

// gets  获取节点列表
func (c *clusterNodeController) nodes(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	nodes, isUpdate, err := c.clusterNodeService.QueryList(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	list := make([]*cluster_dto.ClusterNode, 0, len(nodes))
	for _, node := range nodes {
		list = append(list, &cluster_dto.ClusterNode{
			Name:        node.Name,
			ServiceAddr: node.ServiceAddr,
			AdminAddr:   node.AdminAddr,
			Status:      enum.ClusterNodeStatus(node.Status),
		})
	}

	m := common.Map[string, interface{}]{}
	m["nodes"] = list
	m["is_update"] = isUpdate

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// reset 重置节点信息
func (c *clusterNodeController) reset(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	input := &cluster_dto.ClusterNodeInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.ClusterAddr == "" || input.Source == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("cluster_add is null or source is null"))
		return
	}
	userId := controller.GetUserId(ginCtx)
	if err := c.clusterNodeService.Reset(ginCtx, namespaceId, userId, clusterName, input.ClusterAddr, input.Source); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// post 更新节点信息
func (c *clusterNodeController) put(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	err := c.clusterNodeService.Update(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
