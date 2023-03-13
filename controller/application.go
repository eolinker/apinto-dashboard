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
	"strconv"
)

type applicationController struct {
	applicationService     service.IApplicationService
	applicationAuthService service.IApplicationAuthService
}

func RegisterApplicationRouter(router gin.IRoutes) {
	c := &applicationController{}
	bean.Autowired(&c.applicationService)
	bean.Autowired(&c.applicationAuthService)

	router.GET("/applications", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.lists)
	router.GET("/application/enum", c.lists)
	router.POST("/application", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindApplication), c.createApp)
	router.GET("/application", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.info)
	router.PUT("/application", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindApplication), c.updateApp)
	router.DELETE("/application", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindApplication), c.deleteApp)
	router.GET("/application/onlines", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.onlines)
	router.PUT("/application/online", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindApplication), c.online)
	router.PUT("/application/offline", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypePublish, enum.LogKindApplication), c.offline)
	router.PUT("/application/enable", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindApplication), c.enable)
	router.PUT("/application/disable", GenAccessHandler(access.ApplicationEdit), LogHandler(enum.LogOperateTypeEdit, enum.LogKindApplication), c.disable)
	router.GET("/application/drivers", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.drivers)
	router.GET("/application/auths", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.auths)
	router.GET("/application/auth", GenAccessHandler(access.ApplicationView, access.ApplicationEdit), c.getAuth)
	router.POST("/application/auth", GenAccessHandler(access.ApplicationEdit), c.createAuth)
	router.PUT("/application/auth", GenAccessHandler(access.ApplicationEdit), c.updateAuth)
	router.DELETE("/application/auth", GenAccessHandler(access.ApplicationEdit), c.delAuth)
}

func (a *applicationController) lists(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
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
	userId := GetUserId(ginCtx)
	list, count, err := a.applicationService.AppList(ginCtx, namespaceId, userId, pageNum, pageSize, name)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.ApplicationListOut, 0, len(list))
	for _, application := range list {

		resList = append(resList, &dto.ApplicationListOut{
			Name:       application.Name,
			Id:         application.IdStr,
			Desc:       application.Desc,
			Operator:   application.OperatorName,
			IsDelete:   application.IsDelete,
			UpdateTime: common.TimeToStr(application.UpdateTime),
		})
	}

	data := common.Map[string, interface{}]{}
	data["applications"] = resList
	data["total"] = count
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}
func (a *applicationController) enum(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	list, err := a.applicationService.AppListAll(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.ApplicationEnum, 0, len(list))
	for _, application := range list {

		resList = append(resList, &dto.ApplicationEnum{
			Name: application.Name,
			Id:   application.IdStr,
		})
	}

	data := common.Map[string, interface{}]{}
	data["applications"] = resList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (a *applicationController) info(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	id := ginCtx.Query("app_id")

	info, err := a.applicationService.AppInfoDetails(ginCtx, namespaceId, id)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	customAttrList := make([]dto.ApplicationCustomAttr, 0, len(info.CustomAttr))
	for _, attr := range info.CustomAttr {
		customAttrList = append(customAttrList, dto.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	extraParamList := make([]dto.ApplicationExtraParam, 0, len(info.ExtraParam))
	for _, extra := range info.ExtraParam {
		extraParamList = append(extraParamList, dto.ApplicationExtraParam{
			Key:      extra.Key,
			Value:    extra.Value,
			Conflict: extra.Conflict,
			Position: extra.Position,
		})
	}
	res := dto.ApplicationInfoOut{
		Name:           info.Name,
		Id:             info.IdStr,
		Desc:           info.Desc,
		CustomAttrList: customAttrList,
		ExtraParamList: extraParamList,
	}

	data := common.Map[string, interface{}]{}
	data["application"] = res
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (a *applicationController) createApp(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	input := new(dto.ApplicationInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := input.Check(); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := a.applicationService.CreateApp(ginCtx, namespaceId, userId, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) updateApp(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)

	input := new(dto.ApplicationInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := input.Check(); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := a.applicationService.UpdateApp(ginCtx, namespaceId, userId, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) deleteApp(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	id := ginCtx.Query("app_id")
	userId := GetUserId(ginCtx)
	if err := a.applicationService.DelApp(ginCtx, namespaceId, userId, id); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) onlines(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	id := ginCtx.Query("app_id")
	list, err := a.applicationService.OnlineList(ginCtx, namespaceId, id)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	resp := make([]*dto.OnlineOut, 0, len(list))
	for _, online := range list {
		updateTime := ""
		if !online.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(online.UpdateTime)
		}
		resp = append(resp, &dto.OnlineOut{
			Name:       online.ClusterName,
			Status:     enum.OnlineStatus(online.Status),
			Disable:    online.Disable,
			Env:        online.Env,
			Operator:   online.Operator,
			UpdateTime: updateTime,
		})
	}

	m := common.Map[string, interface{}]{}
	m["clusters"] = resp
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
}

func (a *applicationController) online(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if err := a.applicationService.Online(ginCtx, namespaceId, userId, id, input.ClusterName); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) offline(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if err := a.applicationService.Offline(ginCtx, namespaceId, userId, id, input.ClusterName); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) enable(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if err := a.applicationService.Disable(ginCtx, namespaceId, userId, id, input.ClusterName, false); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) disable(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	id := ginCtx.Query("app_id")
	input := &dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if err := a.applicationService.Disable(ginCtx, namespaceId, userId, id, input.ClusterName, true); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) drivers(ginCtx *gin.Context) {
	driverList := a.applicationAuthService.GetDriversRender()

	drivers := make([]*dto.DriversItem, 0, len(driverList))
	for _, driver := range driverList {
		d := &dto.DriversItem{
			Name:   driver.Name,
			Render: dto.Render(driver.Render),
		}
		drivers = append(drivers, d)
	}
	data := common.Map[string, interface{}]{}
	data["drivers"] = drivers
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (a *applicationController) auths(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	list, err := a.applicationAuthService.GetList(ginCtx, namespaceId, appId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.ApplicationAuthListOut, 0, len(list))

	for _, auth := range list {
		authInfo := &dto.ApplicationAuthListOut{
			Uuid:          auth.Uuid,
			Driver:        auth.Driver,
			ParamPosition: auth.ParamPosition,
			ParamName:     auth.ParamName,
			ParamInfo:     auth.ParamInfo,
			ExpireTime:    auth.ExpireTime,
			Operator:      auth.Operator,
			UpdateTime:    common.TimeToStr(auth.UpdateTime),
			RuleInfo:      auth.RuleInfo,
			IsTransparent: auth.IsTransparent,
		}
		resList = append(resList, authInfo)
	}

	data := common.Map[string, interface{}]{}
	data["auths"] = resList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (a *applicationController) createAuth(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	userId := GetUserId(ginCtx)

	input := &dto.ApplicationAuthInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if input.Position == "" || input.TokenName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("参数位置必填"))
		return
	}
	err := a.applicationAuthService.Create(ginCtx, namespaceId, userId, appId, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) updateAuth(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	uuid := ginCtx.Query("uuid")
	userId := GetUserId(ginCtx)

	input := &dto.ApplicationAuthInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	err := a.applicationAuthService.Update(ginCtx, namespaceId, userId, appId, uuid, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) delAuth(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	userId := GetUserId(ginCtx)

	err := a.applicationAuthService.Delete(ginCtx, namespaceId, userId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (a *applicationController) getAuth(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	appId := ginCtx.Query("app_id")
	uuid := ginCtx.Query("uuid")

	auth, err := a.applicationAuthService.Info(ginCtx, namespaceId, appId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	resAuth := &dto.ApplicationAuthOut{
		Uuid:          auth.Uuid,
		Driver:        auth.Driver,
		ExpireTime:    auth.ExpireTime,
		Operator:      auth.Operator,
		Position:      auth.Position,
		TokenName:     auth.TokenName,
		UpdateTime:    common.TimeToStr(auth.UpdateTime),
		IsTransparent: auth.IsTransparent,
		Config:        dto.AuthConfigProxy(auth.Config),
	}
	data := common.Map[string, interface{}]{}
	data["auth"] = resAuth
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}
