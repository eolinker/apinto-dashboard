package variable_controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//集群下的环境变量

type clusterVariableController struct {
	clusterVariableService variable.IClusterVariableService
}

func RegisterClusterVariableRouter(router gin.IRoutes) {
	c := &clusterVariableController{}
	bean.Autowired(&c.clusterVariableService)

	router.GET("/cluster/:cluster_name/variables", c.gets)
	router.POST("/cluster/:cluster_name/variable", controller.AuditLogHandler(enum.LogOperateTypeCreate, enum.LogKindClusterVariable, c.post))
	router.PUT("/cluster/:cluster_name/variable", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindClusterVariable, c.put))
	router.DELETE("/cluster/:cluster_name/variable", controller.AuditLogHandler(enum.LogOperateTypeDelete, enum.LogKindClusterVariable, c.del))
	router.GET("/cluster/:cluster_name/variable/update-history", c.updateHistory)
	router.POST("/cluster/:cluster_name/variable/sync-conf", c.syncConf)
	router.GET("/cluster/:cluster_name/variable/to-publishs", c.toPublishs)
	router.POST("/cluster/:cluster_name/variable/publish", controller.AuditLogHandler(enum.LogOperateTypePublish, enum.LogKindClusterVariable, c.publish))
	router.GET("/cluster/:cluster_name/variable/publish-history", c.publishHistory)
	router.GET("/cluster/:cluster_name/variable/sync-conf", c.getSyncConf)
}

// gets 获取列表
func (c *clusterVariableController) gets(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)

	variables, err := c.clusterVariableService.GetList(ginCtx, namespaceID, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get ClusterlVariable List fail. err: %s", err.Error())))
		return
	}
	list := make([]*cluster_dto.ClusterVariableItem, 0, len(variables))
	for _, variableInfo := range variables {

		updateTime := ""
		if !variableInfo.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(variableInfo.UpdateTime)
		}
		list = append(list, &cluster_dto.ClusterVariableItem{
			Key:        variableInfo.Key,
			Value:      variableInfo.Value,
			Publish:    enum.PublishType(variableInfo.Publish),
			Desc:       variableInfo.Desc,
			Operator:   variableInfo.Operator,
			UpdateTime: updateTime,
		})
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{"variables": list}))
}

// post 新建
func (c *clusterVariableController) post(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)

	item := &cluster_dto.ClusterVariableItem{}

	if err := ginCtx.BindJSON(item); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	userId := controller.GetUserId(ginCtx)

	if err := c.clusterVariableService.Create(ginCtx, namespaceID, clusterName, userId, item.Key, item.Value, item.Desc); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Create ClusterlVariable fail. err: %s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// put 修改
func (c *clusterVariableController) put(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Alter ClusterVariable fail. err: key can't be nil. ")))
		return
	}

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)

	item := &cluster_dto.ClusterVariableItem{}

	if err := ginCtx.BindJSON(item); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	userId := controller.GetUserId(ginCtx)
	err := c.clusterVariableService.Update(ginCtx, namespaceID, clusterName, userId, key, item.Value)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// del 删除
func (c *clusterVariableController) del(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Delete ClusterVariable fail. err: key can't be nil. ")))
		return
	}
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	err := c.clusterVariableService.Delete(ginCtx, namespaceID, clusterName, userId, key)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Delete ClusterVariable fail. err: %s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// updateHistory 变更历史
func (c *clusterVariableController) updateHistory(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
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

	history, total, err := c.clusterVariableService.QueryHistory(ginCtx, namespaceId, pageNum, pageSize, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	list := make([]*cluster_dto.ClusterHistoryOut, 0, len(history))
	for _, val := range history {
		key := val.OldValue.Key
		if key == "" {
			key = val.NewValue.Key
		}
		list = append(list, &cluster_dto.ClusterHistoryOut{
			Key:        key,
			OldValue:   val.OldValue.Value,
			NewValue:   val.NewValue.Value,
			CreateTime: common.TimeToStr(val.OptTime),
			OptType:    enum.ChangeOptType(val.OptType),
		})
	}

	m := common.Map[string, interface{}]{}
	m["historys"] = list
	m["total"] = total

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// syncConf 同步配置
func (c *clusterVariableController) syncConf(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	conf := new(cluster_dto.SyncConf)
	if err := ginCtx.BindJSON(conf); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if len(conf.Clusters) == 0 || len(conf.Variables) == 0 {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("clusters or variables is null"))
		return
	}
	userId := controller.GetUserId(ginCtx)
	if err := c.clusterVariableService.SyncConf(ginCtx, namespaceId, userId, clusterName, conf); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}

// toPublishs  待发布的环境变量列表
func (c *clusterVariableController) toPublishs(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	background := ginCtx
	list, err := c.clusterVariableService.ToPublishs(background, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	bytes, _ := json.Marshal(list)
	source := common.Base64Encode(bytes)

	toPublishOut := make([]*cluster_dto.VariableToPublishOut, 0, len(list))
	isPublish := true
	for _, publish := range list {
		optType := enum.ChangeOptType(publish.OptType)
		toPublishOut = append(toPublishOut, &cluster_dto.VariableToPublishOut{
			Key:             publish.Key,
			FinishValue:     publish.FinishValue,
			NoReleasedValue: publish.NoReleasedValue,
			CreateTime:      common.TimeToStr(publish.CreateTime),
			OptType:         optType,
		})
	}

	globalVariables, err := c.clusterVariableService.GetList(background, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	m := common.Map[string, interface{}]{}

	defectKeys := make([]string, 0)
	for _, variableInfo := range globalVariables {
		if variableInfo.Publish == 3 {
			defectKeys = append(defectKeys, variableInfo.Key)
			isPublish = false
		}
	}

	if len(defectKeys) > 0 {
		m["unpublish_msg"] = fmt.Sprintf("key为%s的环境变量处于缺失状态不可发布", strings.Join(defectKeys, ","))
	}

	m["variables"] = toPublishOut
	m["is_publish"] = isPublish
	m["source"] = source

	m["version_name"] = time.Now().Format("20060102150405") + "-release"

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))

}

// publish 发布
func (c *clusterVariableController) publish(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	input := &cluster_dto.VariablePublishInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	if input.VersionName == "" || input.Source == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("parameter error"))
		return
	}
	background := ginCtx

	globalVariables, err := c.clusterVariableService.GetList(background, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	defectKeys := make([]string, 0)
	for _, variableInfo := range globalVariables {
		if variableInfo.Publish == 3 {
			defectKeys = append(defectKeys, variableInfo.Key)
		}
	}

	if len(defectKeys) > 0 {
		msg := fmt.Sprintf("key为%s的环境变量处于缺失状态不可发布", strings.Join(defectKeys, ","))
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(msg))
		return
	}

	userId := controller.GetUserId(ginCtx)
	if err = c.clusterVariableService.Publish(background, namespaceId, userId, clusterName, input.VersionName, input.Desc, input.Source); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// publishHistory 发布历史
func (c *clusterVariableController) publishHistory(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
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

	list, total, err := c.clusterVariableService.PublishHistory(ginCtx, namespaceId, pageNum, pageSize, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	historys := make([]*cluster_dto.VariablePublishOut, 0, len(list))
	for _, publish := range list {
		details := make([]*cluster_dto.VariablePublishDetails, 0, len(publish.Details))
		for _, detail := range publish.Details {
			details = append(details, &cluster_dto.VariablePublishDetails{
				Key:        detail.Key,
				OldValue:   detail.OldValue,
				NewValue:   detail.NewValue,
				OptType:    enum.ChangeOptType(detail.OptType),
				CreateTime: common.TimeToStr(detail.CreateTime),
			})
		}
		historys = append(historys, &cluster_dto.VariablePublishOut{
			Id:         publish.Id,
			Name:       publish.Name,
			OptType:    enum.PublishOptType(publish.OptType),
			Operator:   publish.Operator,
			CreateTime: common.TimeToStr(publish.CreateTime),
			Details:    details,
		})
	}

	m := common.Map[string, interface{}]{}
	m["historys"] = historys
	m["total"] = total
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))

}

func (c *clusterVariableController) getSyncConf(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	conf, err := c.clusterVariableService.GetSyncConf(ginCtx, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	syncConf := new(cluster_dto.SyncConf)
	for _, cluster := range conf.Clusters {
		syncConf.Clusters = append(syncConf.Clusters, &cluster_dto.ClusterInput{
			Name: cluster.Name,
			Env:  cluster.Env,
			Id:   cluster.Id,
		})
	}
	for _, variableInfo := range conf.Variables {
		syncConf.Variables = append(syncConf.Variables, &cluster_dto.ClusterVariableSyncConf{
			Id:         variableInfo.Id,
			VariableId: variableInfo.VariableId,
			Key:        variableInfo.Key,
			Value:      variableInfo.Value,
			UpdateTime: common.TimeToStr(variableInfo.UpdateTime),
		})
	}
	m := common.Map[string, interface{}]{}
	m["info"] = syncConf

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
