package open_app_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/openapp"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

type externalApplicationController struct {
	extAppService openapp.IExternalApplicationService
}

func (c externalApplicationController) openApiCheck(ginCtx *gin.Context) {
	//检测openAPI token 并获取相应外部应用的id
	token := ginCtx.GetHeader("Authorization")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	appID, err := c.extAppService.CheckExtAPPToken(ginCtx, namespaceID, token)
	if err != nil {
		ginCtx.Abort()
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("syncAPI fail. err:%s", err))
		return
	}
	ginCtx.Set("appId", appID)
}

var (
	locker             sync.Mutex
	controllerInstance *externalApplicationController
)

func newExternalApplicationController() *externalApplicationController {
	if controllerInstance == nil {
		locker.Lock()
		defer locker.Unlock()
		if controllerInstance == nil {
			controllerInstance = &externalApplicationController{}
			bean.Autowired(&controllerInstance.extAppService)
		}
	}
	return controllerInstance

}
func (e *externalApplicationController) getList(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := e.extAppService.AppList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get external-app list fail. err: %s", err))
		return
	}

	data := common.Map{}
	data["apps"] = list
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (e *externalApplicationController) getInfo(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("id")

	info, err := e.extAppService.AppInfo(ginCtx, namespaceId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get external-app info fail. err: %s", err))
		return
	}

	app := &open_app_dto.ExternalAppInfoOutput{
		Name: info.Name,
		Id:   info.UUID,
		Desc: info.Desc,
	}

	data := common.Map{}
	data["info"] = app
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (e *externalApplicationController) create(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := new(open_app_dto.ExternalAppInfoInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if strings.TrimSpace(input.Name) == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Create external-app fail. err: name can't be null. "))
		return
	}

	if strings.TrimSpace(input.Id) == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Create external-app fail. err: id can't be null. "))
		return
	}

	err := e.extAppService.CreateApp(ginCtx, namespaceId, userId, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Create external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (e *externalApplicationController) edit(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	input := new(open_app_dto.ExternalAppInfoInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Name == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Edit external-app fail. err: name can't be null. "))
		return
	}

	input.Id = uuid
	err := e.extAppService.UpdateApp(ginCtx, namespaceId, userId, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Edit external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (e *externalApplicationController) delete(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.DelApp(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Delete external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (e *externalApplicationController) enable(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.Enable(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Enable external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (e *externalApplicationController) disable(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.Disable(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Disable external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (e *externalApplicationController) flushToken(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	uuid := ginCtx.Query("id")

	err := e.extAppService.FlushToken(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Flush external-app token fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
