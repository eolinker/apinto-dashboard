package controller

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/enum"
	service "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	_ "github.com/eolinker/apinto-dashboard/modules/api/service"
	status_code "github.com/eolinker/apinto-dashboard/modules/api/status-code"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/group/group-dto"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	locker             sync.Mutex
	controllerInstance *apiController
)

type apiController struct {
	apiService service.IAPIService
}

func newApiController() *apiController {
	if controllerInstance == nil {
		locker.Lock()
		defer locker.Unlock()
		if controllerInstance == nil {
			controllerInstance = &apiController{}
			bean.Autowired(&controllerInstance.apiService)
		}
	}
	return controllerInstance

}

func (a *apiController) routerEnum(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	serviceNames := ginCtx.Query("service_names")

	apiList, err := a.apiService.GetAPIListByServiceName(ginCtx, namespaceID, strings.Split(serviceNames, ","))
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	apis := make([]*api_dto.APIEnum, 0, len(apiList))
	for _, item := range apiList {
		api := &api_dto.APIEnum{
			ApiId:   item.UUID,
			APIName: item.Name,
		}

		apis = append(apis, api)
	}

	data := make(map[string]interface{})
	data["apis"] = apis
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (a *apiController) groups(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	parentUUID := ginCtx.Query("parent_uuid")
	queryName := ginCtx.Query("query_name")

	root, apis, err := a.apiService.GetGroups(ginCtx, namespaceID, parentUUID, queryName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resApis := make([]*group_dto.CommonGroupApi, 0, len(apis))
	groups := make([]*group_dto.CommonGroupOut, 0, len(root.CommonGroup))
	for _, group := range root.CommonGroup {
		value := &group_dto.CommonGroupOut{
			UUID:     group.Group.Uuid,
			Name:     group.Group.Name,
			IsDelete: true,
		}
		a.subGroup(value, namespaceID, group.Subgroup)

		for _, subGroup := range value.Children {
			if !subGroup.IsDelete {
				value.IsDelete = false
				break
			}
		}
		if value.IsDelete {
			// 查询一级目录下是否有API
			value.IsDelete = a.apiService.GetAPICountByGroupUUID(context.TODO(), namespaceID, group.Group.Uuid) == 0
		}
		groups = append(groups, value)
	}

	for _, api := range apis {
		resApis = append(resApis, &group_dto.CommonGroupApi{
			Name:      api.Name,
			UUID:      api.UUID,
			Methods:   api.Methods,
			GroupUUID: api.GroupUUID,
		})
	}

	resRoot := &group_dto.CommonGroupRootOut{
		UUID:   root.UUID,
		Name:   root.Name,
		Groups: groups,
	}
	m := make(map[string]interface{})
	m["root"] = resRoot
	m["apis"] = resApis
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

func (a *apiController) subGroup(val *group_dto.CommonGroupOut, namespaceId int, list []*group_model.CommonGroup) {
	if len(list) == 0 {
		return
	}
	for _, group := range list {
		commonGroup := &group_dto.CommonGroupOut{UUID: group.Group.Uuid, Name: group.Group.Name, IsDelete: true}

		a.subGroup(commonGroup, namespaceId, group.Subgroup)
		//若子分组中有不可以删除的分组，则该分组也不能删
		for _, subGroup := range commonGroup.Children {
			if !subGroup.IsDelete {
				commonGroup.IsDelete = false
				if val.IsDelete {
					val.IsDelete = false
				}
				break
			}
		}
		//若子分组均可以删除，则判断该分组下是否有api，若没有，才可以删除
		if commonGroup.IsDelete {
			commonGroup.IsDelete = a.apiService.GetAPICountByGroupUUID(context.TODO(), namespaceId, group.Group.Uuid) == 0
			if !commonGroup.IsDelete && val.IsDelete {
				val.IsDelete = commonGroup.IsDelete
			}
		}

		val.Children = append(val.Children, commonGroup)
	}
}

// routers 获取api列表
func (a *apiController) routers(ginCtx *gin.Context) {

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	groupUUID := ginCtx.Query("group_uuid")
	searchName := ginCtx.Query("search_name")

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

	//获取来源的筛选项， 将 type:id:label 拆分成 ('type',id, 'label')
	sourceIdsStr := ginCtx.Query("source_ids")
	sourceList := strings.Split(sourceIdsStr, ",")
	searchSources := make([]string, 0, len(sourceList))
	for _, source := range sourceList {
		source = strings.ReplaceAll(source, " ", "")
		opts := strings.Split(source, ":")
		if len(opts) != 3 {
			continue
		}
		searchSources = append(searchSources, fmt.Sprintf("('%s',%s,'%s')", opts[0], opts[1], opts[2]))
	}

	apiList, total, err := a.apiService.GetAPIList(ginCtx, namespaceID, groupUUID, searchName, searchSources, pageNum, pageSize)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("GetAPIList fail. err:%s", err.Error()))
		return
	}
	apis := make([]*api_dto.APIListItem, 0, len(apiList))
	for _, item := range apiList {
		publish := make([]*api_dto.APIListItemPublish, 0, len(item.Publish))
		for _, p := range item.Publish {
			publish = append(publish, &api_dto.APIListItemPublish{
				Name:   p.Name,
				Title:  p.Title,
				Status: enum.OnlineStatus(p.Status),
			})
		}
		api := &api_dto.APIListItem{
			GroupUUID:   item.GroupUUID,
			APIUUID:     item.APIUUID,
			APIName:     item.APIName,
			Scheme:      item.Scheme,
			Method:      item.Method,
			ServiceName: item.ServiceName,
			RequestPath: item.RequestPath,
			IsDisable:   item.IsDisable,
			Publish:     publish,
			Source:      item.Source,
			UpdateTime:  common.TimeToStr(item.UpdateTime),
			IsDelete:    item.IsDelete,
		}

		apis = append(apis, api)
	}

	data := make(map[string]interface{})
	data["apis"] = apis
	data["total"] = total
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// getInfo 获取注册中心配置信息
func (a *apiController) getInfo(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取API信息失败 uuid不能为空"))
		return
	}

	info, err := a.apiService.GetAPIVersionInfo(ginCtx, namespaceID, apiUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取API信息失败 err:%s", err.Error()))
		return
	}

	apiInfo := &api_dto.APIInfo{
		ApiName:      info.Api.Name,
		UUID:         info.Api.UUID,
		GroupUUID:    info.Api.GroupUUID,
		Desc:         info.Api.Desc,
		Scheme:       info.Api.Scheme,
		IsDisable:    info.Api.IsDisable,
		RequestPath:  info.Api.RequestPathLabel,
		ServiceName:  info.Version.ServiceName,
		Method:       info.Version.Method,
		ProxyPath:    info.Version.ProxyPath,
		Timeout:      info.Version.Timeout,
		Retry:        info.Version.Retry,
		Hosts:        info.Version.Hosts,
		Match:        info.Version.Match,
		Header:       info.Version.Header,
		TemplateUUID: info.Version.TemplateUUID,
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{"api": apiInfo}))
}

// create 新建
func (a *apiController) create(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := new(api_dto.APIInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, status_code.ApiConfigBindErr, err.Error())
		return
	}

	//API管理器校验参数
	driver := a.apiService.GetAPIDriver(input.Scheme)
	if driver == nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, status_code.ApiSchemeNotExist, fmt.Sprintf("创建API失败, 协议类型不存在"))
		return
	}
	if err := driver.CheckInput(input); err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, status_code.ApiConfigCheckErr, fmt.Sprintf("创建API失败. err:%s", err.Error()))
		return
	}

	uuid, statusCode, err := a.apiService.CreateAPI(ginCtx, namespaceId, userId, input)
	if err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, statusCode, fmt.Sprintf("创建API失败. err:%s", err.Error()))
		return
	}

	data := common.Map{}
	data["uuid"] = uuid
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// alter 修改
func (a *apiController) update(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("更新API失败 err: uuid 不能为空"))
		return
	}

	input := new(api_dto.APIInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//API管理器校验参数
	driver := a.apiService.GetAPIDriver(input.Scheme)
	if driver == nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("更新API失败. 协议无效 "))
		return
	}
	input.UUID = apiUUID
	if err := driver.CheckInput(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("更新API失败 fail. err:%s", err.Error()))
		return
	}

	err := a.apiService.UpdateAPI(ginCtx, namespaceId, userId, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("更新API失败 fail. err:%s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// delete 删除
func (a *apiController) delete(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("DeleteApi fail. err: uuid can't be nil"))
		return
	}

	userId := users.GetUserId(ginCtx)
	err := a.apiService.DeleteAPI(ginCtx, namespaceId, userId, apiUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("DeleteApi fail. err:%s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// batchOnline 批量上线
func (a *apiController) batchOnline(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := &api_dto.ApiBatchInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.OnlineToken == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "online_token can't be nil. ")
		return
	}

	batchOnlineList, err := a.apiService.BatchOnline(ginCtx, namespaceId, userId, input.OnlineToken)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("批量上线失败. err:%s", err.Error()))
		return
	}

	outputList := make([]*api_dto.ApiBatchCheckListItem, 0, len(batchOnlineList))

	for _, item := range batchOnlineList {
		checkItem := &api_dto.ApiBatchCheckListItem{
			ApiName:     item.APIName,
			ClusterName: item.ClusterEnv,
			Status:      item.Status,
			Result:      item.Result,
		}

		outputList = append(outputList, checkItem)
	}

	data := make(map[string]interface{})
	data["list"] = outputList

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// batchOnline 批量下线
func (a *apiController) batchOffline(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := &api_dto.ApiBatchInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if len(input.ApiUUIDs) == 0 || len(input.ClusterNames) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "api_uuids or cluster_names can't be nil. ")
		return
	}

	batchOfflineList, err := a.apiService.BatchOffline(ginCtx, namespaceId, userId, input.ApiUUIDs, input.ClusterNames)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("批量下线API失败. err:%s", err.Error()))
		return
	}

	outputList := make([]*api_dto.ApiBatchCheckListItem, 0, len(batchOfflineList))

	for _, item := range batchOfflineList {
		checkItem := &api_dto.ApiBatchCheckListItem{
			ApiName:     item.APIName,
			ClusterName: item.ClusterEnv,
			Status:      item.Status,
			Result:      item.Result,
		}

		outputList = append(outputList, checkItem)
	}

	data := make(map[string]interface{})
	data["list"] = outputList

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// batchOnline 批量上线检测
func (a *apiController) batchOnlineCheck(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := &api_dto.ApiBatchInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if len(input.ApiUUIDs) == 0 || len(input.ClusterNames) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "api_uuids or cluster_names can't be nil. ")
		return
	}

	batchOnlineList, onlineToken, err := a.apiService.BatchOnlineCheck(ginCtx, namespaceId, userId, input.ApiUUIDs, input.ClusterNames)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("BatchOffline Apis fail. err:%s", err.Error()))
		return
	}

	outputList := make([]*api_dto.ApiBatchOnlineCheckListItem, 0, len(batchOnlineList))

	for _, item := range batchOnlineList {
		checkItem := &api_dto.ApiBatchOnlineCheckListItem{
			ServiceTemplate: item.ServiceTemplate,
			ClusterTitle:    item.ClusterEnv,
			Status:          item.Status,
			Result:          item.Result,
			Solution:        item.Solution,
		}

		outputList = append(outputList, checkItem)
	}

	data := make(map[string]interface{})
	data["list"] = outputList
	if onlineToken != "" {
		data["online_token"] = onlineToken
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// online 上线
func (a *apiController) online(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("uuid can't be nil"))
		return
	}
	input := &api_dto.PublishInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	router, err := a.apiService.OnlineAPI(ginCtx, namespaceId, userId, apiUUID, input.ClusterNames)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	} else if len(router) == 0 {
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
		return
	}

	msg := ""
	if err != nil {
		msg = err.Error()
	}

	m := make(map[string]interface{})
	m["router"] = router
	ginCtx.JSON(http.StatusOK, controller.Result{
		Code: -1,
		Data: m,
		Msg:  msg,
	})

}

// online 下线
func (a *apiController) offline(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("uuid 不可以为空"))
		return
	}
	input := &api_dto.PublishInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if _, err := a.apiService.OfflineAPI(ginCtx, namespaceId, userId, apiUUID, input.ClusterNames); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *apiController) getOnlineInfo(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	apiUUID := ginCtx.Query("uuid")
	if apiUUID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取发布管理信息失败: uuid 不能为空"))
		return
	}
	apiInfo, clustersPublish, err := a.apiService.OnlineInfo(ginCtx, namespaceId, apiUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取发布管理信息失败: %s", err.Error()))
		return
	}

	info := &api_dto.ApiPublishInfo{
		Name:      apiInfo.Api.Name,
		ID:        apiInfo.Api.UUID,
		Scheme:    apiInfo.Api.Scheme,
		Method:    apiInfo.Version.Method,
		Path:      apiInfo.Api.RequestPathLabel,
		Service:   apiInfo.Version.ServiceName,
		ProxyPath: apiInfo.Version.ProxyPath,
		Desc:      apiInfo.Api.Desc,
	}

	clusters := make([]*api_dto.ApiPublishCluster, 0, len(clustersPublish))
	for _, clu := range clustersPublish {
		clusters = append(clusters, &api_dto.ApiPublishCluster{
			Name:       clu.ClusterName,
			Env:        clu.ClusterEnv,
			Title:      clu.ClusterTitle,
			Status:     enum.OnlineStatus(clu.Status),
			Operator:   clu.Operator,
			UpdateTime: clu.UpdateTime,
		})
	}
	m := make(map[string]interface{})
	m["info"] = info
	m["clusters"] = clusters
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// getOnlineList 上线管理列表
func (a *apiController) getOnlineList(ginCtx *gin.Context) {
	//namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	//apiUUID := ginCtx.Query("uuid")
	//if apiUUID == "" {
	//	controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("DeleteApi fail. err: uuid can't be nil"))
	//	return
	//}
	//
	//list, err := a.apiService.OnlineList(ginCtx, namespaceId, apiUUID)
	//if err != nil {
	//	controller.ErrorJson(ginCtx, http.StatusOK, err.Logger())
	//	return
	//}
	//
	//resp := make([]*online_dto.OnlineOut, 0, len(list))
	//for _, item := range list {
	//	updateTime := ""
	//	if !item.UpdateTime.IsZero() {
	//		updateTime = common.TimeToStr(item.UpdateTime)
	//	}
	//	resp = append(resp, &online_dto.OnlineOut{
	//		Name:       item.ClusterName,
	//		Env:        item.ClusterEnv,
	//		Status:     enum.OnlineStatus(item.Status),
	//		Disable:    item.Disable,
	//		Operator:   item.Operator,
	//		UpdateTime: updateTime,
	//	})
	//}
	//
	//m := make(map[string]interface{})
	//m["clusters"] = resp
	//ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// getSourceList 获取来源列表
func (a *apiController) getSourceList(ginCtx *gin.Context) {
	items, err := a.apiService.GetSource(ginCtx)
	if err != nil {
		log.Error(err)
	}
	data := make(map[string]interface{})
	data["list"] = items
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// getImportCheckList 获取导入Swagger后的检查列表
func (a *apiController) getImportCheckList(ginCtx *gin.Context) {
	fileInfo, err := ginCtx.FormFile("file")
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger get file fail. err: %s. ", err))
		return
	}
	file, err := fileInfo.Open()
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger open file fail. err: %s. ", err))
		return
	}
	fileData := make([]byte, fileInfo.Size)
	_, err = file.Read(fileData)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger  read file fail. err: %s. ", err))
		return
	}
	defer file.Close()

	//TODO format 以后待处理
	//TODO is_base64 以后待处理

	groupID := ginCtx.PostForm("group")
	if groupID == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger fail. err: group can't be null. "))
		return
	}
	serviceName := ginCtx.PostForm("upstream")
	if serviceName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger fail. err: upstream can't be null. "))
		return
	}
	requestPrefix := ginCtx.PostForm("request_prefix")

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	checkList, token, err := a.apiService.GetImportCheckList(ginCtx, namespaceID, fileData, groupID, serviceName, requestPrefix)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("import swagger fail. err: %s. ", err))
		return
	}

	resList := make([]api_dto.ImportAPIListItem, 0)
	for _, item := range checkList {
		resList = append(resList, api_dto.ImportAPIListItem{
			Id:     item.Id,
			Name:   item.Name,
			Method: item.Method,
			Path:   item.Path,
			Desc:   item.Desc,
			Status: enum.ImportStatusType(item.Status),
		})
	}

	data := make(map[string]interface{})
	data["apis"] = resList
	data["token"] = token
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))

}

// importAPI 导入Swagger文档
func (a *apiController) importAPI(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	input := new(api_dto.ImportAPIInfos)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	err := a.apiService.ImportAPI(ginCtx, namespaceID, userId, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("importAPI fail. err: %s. ", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (a *apiController) checkApiExist(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	apiInfo, err := a.apiService.GetAPIInfo(ginCtx, namespaceID, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("CheckApiExist fail. err: %s. ", err))
		return
	}

	isExist := false
	if apiInfo != nil {
		isExist = true
	}
	data := common.Map{}
	data["is_exist"] = isExist
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
