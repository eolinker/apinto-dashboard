package open_api_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	RegisterRouter(newClusterOpenController())
}

func newClusterOpenController() *clusterOpenController {
	c := &clusterOpenController{}
	bean.Autowired(&c.clusterService)
	return c
}

type clusterOpenController struct {
	clusterService cluster.IClusterService
}

func (c *clusterOpenController) Name() string {
	return "cluster"
}

func (c *clusterOpenController) Apis() []pm3.Api {
	return []pm3.Api{
		{
			Method:      http.MethodGet,
			Path:        "/api2/clusters/simple",
			HandlerFunc: c.simpleClusters,
			Authority:   pm3.Public,
		},
	}
}

func (c *clusterOpenController) simpleClusters(ginCtx *gin.Context) {
	clusters, err := c.clusterService.SimpleCluster(ginCtx, namespace_controller.GetNamespaceId(ginCtx))
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["clusters"] = clusters

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
