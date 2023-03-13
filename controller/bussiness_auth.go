package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type bussinessAuthController struct {
	bussinessAuthService service.IBussinessAuthService
}

func RegisterBussinessAuthRouter(router gin.IRouter) {
	b := &bussinessAuthController{}
	bean.Autowired(&b.bussinessAuthService)

	router.GET("/_system/activation/info", GenAccessHandler(access.AuthorizationView, access.AuthorizationEdit), b.acivationInfo)
	router.GET("/_system/run-info", b.getRunInfo)
	router.GET("/_system/mac", b.getMac)
	router.POST("/_system/activation", b.activation)
	router.POST("/_system/reactivation", GenAccessHandler(access.AuthorizationEdit), b.reactivation)
}

func (b *bussinessAuthController) acivationInfo(ginCtx *gin.Context) {
	info, err := b.bussinessAuthService.GetActivationInfo(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("get activationInfi fail. err: %s ", err)))
		return
	}

	data := make(map[string]interface{})
	data["infos"] = formatActivationInfo(info)
	data["title"] = enum.EditionTitle(info.Edition)
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))

}

func (b *bussinessAuthController) getRunInfo(ginCtx *gin.Context) {
	runInfo, err := b.bussinessAuthService.GetActivationInfo(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("get runInfo fail. err: %s ", err)))
		return
	}

	data := make(map[string]interface{})
	data["info"] = &dto.SystemRunInfo{
		Title:       enum.EditionTitle(runInfo.Edition),
		DashboardID: runInfo.DashboardID,
		Version:     version.GetVersion(),
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))

}

func (b *bussinessAuthController) getMac(ginCtx *gin.Context) {
	mac, err := b.bussinessAuthService.GetMachineCode(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("get machine_code fail. err: %s. ", err)))
		return
	}

	data := make(map[string]interface{})
	info := make(map[string]interface{})
	info["mac"] = mac
	data["info"] = info
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (b *bussinessAuthController) activation(ginCtx *gin.Context) {
	fileInfo, err := ginCtx.FormFile("auth_file")
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation get auth_file fail. err: %s. ", err)))
		return
	}
	file, err := fileInfo.Open()
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation open auth_file fail. err: %s. ", err)))
		return
	}
	certData := make([]byte, fileInfo.Size)
	_, err = file.Read(certData)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation read auth_file fail. err: %s. ", err)))
		return
	}
	defer file.Close()

	info, err := b.bussinessAuthService.ActivateCert(ginCtx, certData)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation fail. err: %s. ", err)))
		return
	}

	data := make(map[string]interface{})
	data["infos"] = formatActivationInfo(info)
	data["title"] = enum.EditionTitle(info.Edition)
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (b *bussinessAuthController) reactivation(ginCtx *gin.Context) {
	fileInfo, err := ginCtx.FormFile("auth_file")
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation get auth_file fail. err: %s. ", err)))
		return
	}
	file, err := fileInfo.Open()
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation open auth_file fail. err: %s. ", err)))
		return
	}
	certData := make([]byte, fileInfo.Size)
	_, err = file.Read(certData)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("activation read auth_file fail. err: %s. ", err)))
		return
	}
	defer file.Close()

	info, err := b.bussinessAuthService.ReActivateCert(ginCtx, certData)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("re-activation fail. err: %s. ", err)))
		return
	}

	data := make(map[string]interface{})
	data["infos"] = formatActivationInfo(info)
	data["title"] = enum.EditionTitle(info.Edition)
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func formatActivationInfo(info *model.ActivationInfo) []*dto.ActivationInfoItem {
	items := make([]*dto.ActivationInfoItem, 0, 5)

	items = append(items, &dto.ActivationInfoItem{
		Key:   "授权对象",
		Value: info.Company,
	})

	items = append(items, &dto.ActivationInfoItem{
		Key:   "节点数量",
		Value: formatNodeTitle(info.NodeCount),
	})
	items = append(items, &dto.ActivationInfoItem{
		Key:   "控制台数量",
		Value: formatNodeTitle(info.ControllerCount),
	})
	items = append(items, &dto.ActivationInfoItem{
		Key:   "控制台ID",
		Value: info.DashboardID,
	})
	validTime := time.Unix(info.EndTime, 0).Format("2006-01-02")
	items = append(items, &dto.ActivationInfoItem{
		Key:   "有效期",
		Value: validTime,
	})

	return items
}

func formatNodeTitle(count int) string {
	if count == -1 {
		return "不限"
	}
	return strconv.Itoa(count)
}
