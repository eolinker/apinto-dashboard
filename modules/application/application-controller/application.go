package application_controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"sync"

	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/application/application-dto"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/online/online-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	locker             sync.Mutex
	controllerInstance *applicationController
)

type applicationController struct {
	applicationService application.IApplicationService
}

func newApplicationController() *applicationController {
	if controllerInstance == nil {
		locker.Lock()
		defer locker.Unlock()
		if controllerInstance == nil {
			controllerInstance = &applicationController{}
			bean.Autowired(&controllerInstance.applicationService)
		}
	}
	return controllerInstance
}

func (a *applicationController) lists(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	pageNumStr := ginCtx.Query("page_num")
	pageSizeStr := ginCtx.Query("page_size")
	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	name := ginCtx.Query("name")
	userId := users.GetUserId(ginCtx)

	clustersStr := ginCtx.Query("clusters")
	clusters := make([]string, 0)
	if clustersStr != "" {
		err := json.Unmarshal([]byte(clustersStr), &clusters)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
	}

	list, count, err := a.applicationService.AppList(ginCtx, namespaceId, userId, pageNum, pageSize, name, clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resList := make([]*application_dto.ApplicationListItem, 0, len(list))
	for _, item := range list {
		publish := make([]*application_dto.APPListItemPublish, 0, len(item.Publish))
		for _, p := range item.Publish {
			publish = append(publish, &application_dto.APPListItemPublish{
				Name:   p.Name,
				Title:  p.Title,
				Status: enum.OnlineStatus(p.Status),
			})
		}
		resList = append(resList, &application_dto.ApplicationListItem{
			Name:       item.Name,
			Id:         item.Uuid,
			Desc:       item.Desc,
			Publish:    publish,
			Operator:   item.OperatorName,
			UpdateTime: common.TimeToStr(item.UpdateTime),
			IsDelete:   item.IsDelete,
		})
	}

	data := common.Map{}
	data["applications"] = resList
	data["total"] = count
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) enum(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := a.applicationService.AppEnumList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resList := make([]*application_dto.ApplicationEnum, 0, len(list))
	for _, applicationInfo := range list {

		resList = append(resList, &application_dto.ApplicationEnum{
			Title: applicationInfo.Name,
			Id:    applicationInfo.Uuid,
		})
	}

	data := common.Map{}
	data["applications"] = resList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) info(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	id := ginCtx.Query("app_id")

	info, err := a.applicationService.AppInfoDetails(ginCtx, namespaceId, id)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	customAttrList := make([]application_dto.ApplicationCustomAttr, 0, len(info.CustomAttr))
	paramList := make([]application_dto.ExtraParam, 0, len(info.Params))
	for _, attr := range info.CustomAttr {
		customAttrList = append(customAttrList, application_dto.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	for _, param := range info.Params {
		paramList = append(paramList, application_dto.ExtraParam{
			Key:      param.Key,
			Value:    param.Value,
			Conflict: param.Conflict,
			Position: param.Position,
		})
	}

	res := application_dto.ApplicationInfoOut{
		Name:           info.Name,
		Id:             info.Uuid,
		Desc:           info.Desc,
		CustomAttrList: customAttrList,
		Params:         paramList,
	}

	data := common.Map{}
	data["application"] = res
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) createApp(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	input := new(application_dto.ApplicationInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := input.Check(); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := a.applicationService.CreateApp(ginCtx, namespaceId, userId, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) updateApp(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := new(application_dto.ApplicationInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := input.Check(); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := a.applicationService.UpdateApp(ginCtx, namespaceId, userId, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) deleteApp(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	id := ginCtx.Query("app_id")
	userId := users.GetUserId(ginCtx)
	if err := a.applicationService.DelApp(ginCtx, namespaceId, userId, id); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) online(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &online_dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if err := a.applicationService.Online(ginCtx, namespaceId, userId, id, input.ClusterNames); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) offline(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &online_dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if err := a.applicationService.Offline(ginCtx, namespaceId, userId, id, input.ClusterNames); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) getOnlineInfo(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取发布管理信息失败: uuid 不能为空"))
		return
	}
	appInfo, clustersPublish, err := a.applicationService.OnlineInfo(ginCtx, namespaceId, apiUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取发布管理信息失败: %s", err.Error()))
		return
	}

	info := &application_dto.AppPublishInfo{
		Name: appInfo.Name,
		ID:   appInfo.Uuid,
		Desc: appInfo.Desc,
	}

	clusters := make([]*application_dto.AppPublishCluster, 0, len(clustersPublish))
	for _, clu := range clustersPublish {
		clusters = append(clusters, &application_dto.AppPublishCluster{
			Name:       clu.Name,
			Env:        clu.Env,
			Title:      clu.Title,
			Status:     enum.OnlineStatus(clu.Status),
			Operator:   clu.Updater,
			UpdateTime: clu.UpdateTime,
		})
	}
	m := make(map[string]interface{})
	m["info"] = info
	m["clusters"] = clusters
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

func (a *applicationController) drivers(ginCtx *gin.Context) {
	driverList := a.applicationService.GetDriversRender()

	drivers := make([]*application_dto.AuthDriversItem, 0, len(driverList))
	for _, driver := range driverList {
		d := &application_dto.AuthDriversItem{
			Name:   driver.Name,
			Render: application_dto.Render(driver.Render),
		}
		drivers = append(drivers, d)
	}
	data := common.Map{}
	data["drivers"] = drivers
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) auths(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	list, err := a.applicationService.GetAuthList(ginCtx, namespaceId, appId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resList := make([]*application_dto.ApplicationAuthListOut, 0, len(list))

	for _, auth := range list {
		authInfo := &application_dto.ApplicationAuthListOut{
			Uuid:           auth.UUID,
			Title:          auth.Title,
			Driver:         auth.Driver,
			HideCredential: auth.HideCredential,
			ExpireTime:     auth.ExpireTime,
			Operator:       auth.Operator,
			UpdateTime:     common.TimeToStr(auth.UpdateTime),
		}
		resList = append(resList, authInfo)
	}

	data := common.Map{}
	data["auths"] = resList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) createAuth(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	userId := users.GetUserId(ginCtx)

	input := &application_dto.ApplicationAuthInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if input.Position == "" || input.TokenName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "参数位置必填")
		return
	}
	err := a.applicationService.CreateAuth(ginCtx, namespaceId, userId, appId, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) updateAuth(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	uuid := ginCtx.Query("uuid")
	userId := users.GetUserId(ginCtx)

	input := &application_dto.ApplicationAuthInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	err := a.applicationService.UpdateAuth(ginCtx, namespaceId, userId, appId, uuid, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) delAuth(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	userId := users.GetUserId(ginCtx)

	err := a.applicationService.DeleteAuth(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *applicationController) getAuth(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	uuid := ginCtx.Query("uuid")

	auth, err := a.applicationService.AuthInfo(ginCtx, namespaceId, appId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	resAuth := &application_dto.ApplicationAuthOut{
		Title:          auth.Title,
		Uuid:           auth.Uuid,
		Driver:         auth.Driver,
		ExpireTime:     auth.ExpireTime,
		Operator:       auth.Operator,
		Position:       auth.Position,
		TokenName:      auth.TokenName,
		UpdateTime:     common.TimeToStr(auth.UpdateTime),
		HideCredential: auth.IsTransparent,
		Config:         application_dto.AuthConfigProxy(auth.Config),
	}
	data := common.Map{}
	data["auth"] = resAuth
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *applicationController) getAuthDetails(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	uuid := ginCtx.Query("uuid")

	result, err := a.applicationService.AuthDetails(ginCtx, namespaceId, appId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	data := common.Map{}
	data["details"] = result
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
