package cluster_controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type clusterCertificateController struct {
	clusterCertificateService cluster.IClusterCertificateService
}

func newClusterCertificateController() *clusterCertificateController {
	c := &clusterCertificateController{}
	bean.Autowired(&c.clusterCertificateService)
	return c
}

// gets 获取证书列表
func (c *clusterCertificateController) gets(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	list, err := c.clusterCertificateService.QueryList(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	dtoList := make([]*cluster_dto.ClusterCertificateOut, 0, len(list))
	for _, val := range list {
		cert, err := common.ParseCert(val.Key, val.Pem)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		dtoList = append(dtoList, &cluster_dto.ClusterCertificateOut{
			Id:           val.Id,
			ClusterId:    val.ClusterId,
			Name:         cert.Leaf.Subject.CommonName,
			ValidTime:    common.TimeToStr(cert.Leaf.NotAfter),
			OperatorName: val.OperatorName,
			CreateTime:   common.TimeToStr(val.CreateTime),
			UpdateTime:   common.TimeToStr(val.UpdateTime),
		})
	}
	m := common.Map[string, interface{}]{}
	m["certificates"] = dtoList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// post 新增
func (c *clusterCertificateController) post(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	operator := users.GetUserId(ginCtx)
	input := &cluster_dto.ClusterCertificateInput{}
	err := ginCtx.BindJSON(input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if len(input.Key) == 0 || len(input.Pem) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "key or pem is null")
		return
	}

	pem, _ := common.Base64Decode(input.Pem)
	key, _ := common.Base64Decode(input.Key)

	if err = c.clusterCertificateService.Insert(ginCtx, operator, namespaceId, clusterName, string(key), string(pem)); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// put 修改
func (c *clusterCertificateController) put(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	certificateIdStr := ginCtx.Param("certificate_id")
	certificateId, _ := strconv.Atoi(certificateIdStr)
	operator := users.GetUserId(ginCtx)
	if certificateId <= 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "certificate_id is 0")
		return
	}
	input := &cluster_dto.ClusterCertificateInput{}
	err := ginCtx.BindJSON(input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if len(input.Key) == 0 || len(input.Pem) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "key or pem is null")
		return
	}

	pem, _ := common.Base64Decode(input.Pem)
	key, _ := common.Base64Decode(input.Key)

	if err = c.clusterCertificateService.Update(ginCtx, operator, namespaceId, certificateId, clusterName, string(key), string(pem)); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// del 删除
func (c *clusterCertificateController) del(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	certificateIdStr := ginCtx.Param("certificate_id")
	certificateId, _ := strconv.Atoi(certificateIdStr)

	err := c.clusterCertificateService.DeleteById(ginCtx, namespaceId, clusterName, certificateId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}
