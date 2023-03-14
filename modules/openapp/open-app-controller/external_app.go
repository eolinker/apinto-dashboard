package open_app_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/openapp"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type externalApplicationController struct {
	extAppService openapp.IExternalApplicationService
}

func RegisterExternalApplicationRouter(router gin.IRoutes) {
	e := &externalApplicationController{}
	bean.Autowired(&e.extAppService)

	router.GET("/external-apps", controller.GenAccessHandler(access.ExtAPPEdit, access.ExtAPPView), e.getList)
	router.GET("/external-app", controller.GenAccessHandler(access.ExtAPPEdit, access.ExtAPPView), e.getInfo)
	router.POST("/external-app", controller.GenAccessHandler(access.ExtAPPEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindExtAPP), e.create)
	router.PUT("/external-app", controller.GenAccessHandler(access.ExtAPPEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindExtAPP), e.edit)
	router.DELETE("/external-app", controller.GenAccessHandler(access.ExtAPPEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindExtAPP), e.delete)
	router.PUT("/external-app/enable", controller.GenAccessHandler(access.ExtAPPEdit), e.enable)
	router.PUT("/external-app/disable", controller.GenAccessHandler(access.ExtAPPEdit), e.disable)
	router.PUT("/external-app/token", controller.GenAccessHandler(access.ExtAPPEdit), e.flushToken)
}

func (e *externalApplicationController) getList(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := e.extAppService.AppList(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Get external-app list fail. err: %s", err)))
		return
	}

	data := common.Map[string, interface{}]{}
	data["apps"] = list
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (e *externalApplicationController) getInfo(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("id")

	info, err := e.extAppService.AppInfo(ginCtx, namespaceId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Get external-app info fail. err: %s", err)))
		return
	}

	app := &open_app_dto.ExternalAppInfoOutput{
		Name: info.Name,
		Id:   info.UUID,
		Desc: info.Desc,
	}

	data := common.Map[string, interface{}]{}
	data["info"] = app
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (e *externalApplicationController) create(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	input := new(open_app_dto.ExternalAppInfoInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if strings.TrimSpace(input.Name) == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Create external-app fail. err: name can't be null. ")))
		return
	}

	if strings.TrimSpace(input.Id) == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Create external-app fail. err: id can't be null. ")))
		return
	}

	err := e.extAppService.CreateApp(ginCtx, namespaceId, userId, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Create external-app fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (e *externalApplicationController) edit(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	input := new(open_app_dto.ExternalAppInfoInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if input.Name == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Edit external-app fail. err: name can't be null. ")))
		return
	}

	input.Id = uuid
	err := e.extAppService.UpdateApp(ginCtx, namespaceId, userId, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Edit external-app fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (e *externalApplicationController) delete(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.DelApp(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Delete external-app fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (e *externalApplicationController) enable(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.Enable(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Enable external-app fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (e *externalApplicationController) disable(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.Disable(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Disable external-app fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (e *externalApplicationController) flushToken(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.FlushToken(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Flush external-app token fail. err: %s", err)))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}
