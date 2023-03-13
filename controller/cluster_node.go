package controller

import (
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type clusterNodeController struct {
	clusterNodeService service.IClusterNodeService
}

func RegisterClusterNodeRouter(router gin.IRoutes) {
	c := &clusterNodeController{}
	bean.Autowired(&c.clusterNodeService)

	router.GET("/cluster/:cluster_name/nodes", genAccessHandler(access.ClusterView, access.ClusterEdit), c.nodes)
	router.POST("/cluster/:cluster_name/node/reset", genAccessHandler(access.ClusterEdit), c.reset)
	router.PUT("/cluster/:cluster_name/node", genAccessHandler(access.ClusterEdit), c.put)
}

// gets  获取节点列表
func (c *clusterNodeController) nodes(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	nodes, isUpdate, err := c.clusterNodeService.QueryList(ginCtx, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	list := make([]*dto.ClusterNode, 0, len(nodes))
	for _, node := range nodes {
		list = append(list, &dto.ClusterNode{
			Name:        node.Name,
			ServiceAddr: node.ServiceAddr,
			AdminAddr:   node.AdminAddr,
			Status:      enum.ClusterNodeStatus(node.Status),
		})
	}

	m := common.Map[string, interface{}]{}
	m["nodes"] = list
	m["is_update"] = isUpdate

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
}

// reset 重置节点信息
func (c *clusterNodeController) reset(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	input := &dto.ClusterNodeInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if input.ClusterAddr == "" || input.Source == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("cluster_add is null or source is null"))
		return
	}
	userId := getUserId(ginCtx)
	if err := c.clusterNodeService.Reset(ginCtx, namespaceId, userId, clusterName, input.ClusterAddr, input.Source); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// post 更新节点信息
func (c *clusterNodeController) put(ginCtx *gin.Context) {

	namespaceId := getNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	err := c.clusterNodeService.Update(ginCtx, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}
