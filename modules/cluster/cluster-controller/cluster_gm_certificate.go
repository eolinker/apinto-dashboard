package cluster_controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type gmCertificateController struct {
	clusterCertificateService cluster.IGmCertificateService
}

func newGmCertificateController() *gmCertificateController {
	c := &gmCertificateController{}
	bean.Autowired(&c.clusterCertificateService)
	return c
}

func parseCert(key string, cert string) (string, *cluster_dto.CertificateInfo, error) {
	info, err := common.ParseGMCert(key, cert)
	if err != nil {
		return "", nil, err
	}
	dnsNames := info.Leaf.DNSNames
	if dnsNames == nil && info.Leaf.IPAddresses != nil {
		dnsNames = make([]string, 0, len(info.Leaf.IPAddresses))
		for _, ip := range info.Leaf.IPAddresses {
			dnsNames = append(dnsNames, ip.String())
		}
	}
	if dnsNames == nil {
		dnsNames = make([]string, 0)
		dnsNames = append(dnsNames, info.Leaf.Subject.CommonName)
	}
	return info.Leaf.Subject.CommonName, &cluster_dto.CertificateInfo{
		DNSName:   dnsNames,
		ValidTime: common.TimeToStr(info.Leaf.NotAfter),
	}, nil

}

// gets 获取证书列表
func (c *gmCertificateController) gets(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	list, err := c.clusterCertificateService.QueryList(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	dtoList := make([]*cluster_dto.GMCertificateItem, 0, len(list))
	for _, val := range list {
		commonName, signCert, err := parseCert(val.SignKey, val.SignCert)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		_, encCert, err := parseCert(val.EncKey, val.EncCert)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		dtoList = append(dtoList, &cluster_dto.GMCertificateItem{
			Id:           val.Id,
			ClusterId:    val.ClusterId,
			Name:         commonName,
			SignCert:     signCert,
			EncCert:      encCert,
			OperatorName: val.OperatorName,
			CreateTime:   common.TimeToStr(val.CreateTime),
			UpdateTime:   common.TimeToStr(val.UpdateTime),
		})
	}
	m := common.Map{}
	m["certificates"] = dtoList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// gets 获取证书列表
func (c *gmCertificateController) get(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	certificateIdStr := ginCtx.Param("certificate_id")
	certificateId, _ := strconv.Atoi(certificateIdStr)
	info, err := c.clusterCertificateService.Info(ginCtx, namespaceId, certificateId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"certificate": info,
	}))
}

// post 新增
func (c *gmCertificateController) post(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	operator := users.GetUserId(ginCtx)
	input := &cluster_dto.GMCertificateInput{}
	err := ginCtx.BindJSON(input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if len(input.SignKey) == 0 || len(input.SignCert) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "sign key or sign cert is null")
		return
	}

	if len(input.EncKey) == 0 || len(input.EncCert) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "enc key or enc cert is null")
		return
	}

	signKey, _ := common.Base64Decode(input.SignKey)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("sign key is invalid:%s", err.Error()))
		return
	}
	signCert, err := common.Base64Decode(input.SignCert)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("sign cert is invalid:%s", err.Error()))
		return
	}

	encKey, err := common.Base64Decode(input.EncKey)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("enc key is invalid:%s", err.Error()))
		return
	}
	encCert, err := common.Base64Decode(input.EncCert)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("enc cert is invalid:%s", err.Error()))
		return
	}

	if err = c.clusterCertificateService.Insert(ginCtx, operator, namespaceId, clusterName, string(signKey), string(signCert), string(encKey), string(encCert)); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// put 修改
func (c *gmCertificateController) put(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	certificateIdStr := ginCtx.Param("certificate_id")
	certificateId, _ := strconv.Atoi(certificateIdStr)
	operator := users.GetUserId(ginCtx)
	if certificateId <= 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "certificate_id is 0")
		return
	}
	input := &cluster_dto.GMCertificateInput{}
	err := ginCtx.BindJSON(input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if len(input.SignKey) == 0 || len(input.SignCert) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "sign key or sign cert is null")
		return
	}

	if len(input.EncKey) == 0 || len(input.EncCert) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "enc key or enc cert is null")
		return
	}

	signKey, _ := common.Base64Decode(input.SignKey)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("sign key is invalid:%s", err.Error()))
		return
	}
	signCert, err := common.Base64Decode(input.SignCert)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("sign cert is invalid:%s", err.Error()))
		return
	}

	encKey, err := common.Base64Decode(input.EncKey)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("enc key is invalid:%s", err.Error()))
		return
	}
	encCert, err := common.Base64Decode(input.EncCert)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("enc cert is invalid:%s", err.Error()))
		return
	}

	if err = c.clusterCertificateService.Update(ginCtx, operator, namespaceId, certificateId, clusterName, string(signKey), string(signCert), string(encKey), string(encCert)); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// del 删除
func (c *gmCertificateController) del(ginCtx *gin.Context) {
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
