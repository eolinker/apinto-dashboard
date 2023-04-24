package controller

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/api"
	api_model "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	cluster_model "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	monitor_dto "github.com/eolinker/apinto-dashboard/modules/monitor/monitor-dto"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"

	"github.com/eolinker/apinto-dashboard/modules/monitor"
	"github.com/eolinker/apinto-dashboard/modules/monitor/model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

type monitorController struct {
	monitorService monitor.IMonitorService
	monStatistics  monitor.IMonitorStatistics
	apiService     api.IAPIService
	appService     application.IApplicationService
	clusterService cluster.IClusterService
	service        upstream.IService
}

func RegisterMonitorRouter(router gin.IRoutes) {
	m := &monitorController{}
	bean.Autowired(&m.monitorService)
	bean.Autowired(&m.monStatistics)
	bean.Autowired(&m.apiService)
	bean.Autowired(&m.appService)
	bean.Autowired(&m.service)
	bean.Autowired(&m.clusterService)

	prefix := "/monitor"

	router.GET("/monitor/partitions", m.getPartitionList)
	router.GET("/monitor/partition", m.getPartitionInfo)
	router.POST("/monitor/partition", m.createPartition)
	router.PUT("/monitor/partition", m.editPartition)
	router.DELETE("/monitor/partition", m.delPartition)
	//总览
	router.POST(prefix+"/overview/summary", gzip.Gzip(gzip.DefaultCompression), m.overviewSummary)
	router.POST(prefix+"/overview/invoke", gzip.Gzip(gzip.DefaultCompression), m.overviewInvoke)
	router.POST(prefix+"/overview/message", gzip.Gzip(gzip.DefaultCompression), m.overviewMessage)
	router.POST(prefix+"/overview/top", gzip.Gzip(gzip.DefaultCompression), m.overviewTop)

	//数据统计
	router.POST(prefix+"/:dataType", gzip.Gzip(gzip.DefaultCompression), m.getStatistics)
	//数据统计详情
	router.POST(prefix+"/:dataType/details", gzip.Gzip(gzip.DefaultCompression), m.getStatisticsDetails)
	//数据统计详情-数据
	router.POST(prefix+"/:dataType/details/:detailsType", gzip.Gzip(gzip.DefaultCompression), m.getStatistics)
	//数据统计详情-数据-趋势图
	router.POST(prefix+"/:dataType/details/:detailsType/trend", gzip.Gzip(gzip.DefaultCompression), m.getStatisticsDetails)
}

func (m *monitorController) getPartitionList(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := m.monitorService.PartitionList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get monitor partition list fail. err: %s", err))
		return
	}

	data := common.Map[string, interface{}]{}
	data["partitions"] = list
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (m *monitorController) getPartitionInfo(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")

	info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get monitor partition info fail. err: %s", err))
		return
	}

	respInfo := &monitor_dto.MonitorPartitionInfoProxy{
		Name:         info.Name,
		SourceType:   info.SourceType,
		Config:       info.Config,
		Env:          info.Env,
		ClusterNames: info.ClusterNames,
	}

	data := common.Map[string, interface{}]{}
	data["info"] = respInfo
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (m *monitorController) createPartition(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	inputProxy := new(monitor_dto.MonitorPartitionInfoProxy)
	if err := ginCtx.BindJSON(inputProxy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	var err error
	inputProxy.Config, err = m.monitorService.CheckInput(inputProxy.SourceType, inputProxy.Config)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	info, err := m.monitorService.CreatePartition(ginCtx, namespaceId, userId, inputProxy)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Create monitor partition fail. err: %s", err))
		return
	}

	data := common.Map[string, interface{}]{}
	data["info"] = info
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (m *monitorController) editPartition(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("uuid")

	inputProxy := new(monitor_dto.MonitorPartitionInfoProxy)
	if err := ginCtx.BindJSON(inputProxy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	var err error
	inputProxy.Config, err = m.monitorService.CheckInput(inputProxy.SourceType, inputProxy.Config)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	info, err := m.monitorService.UpdatePartition(ginCtx, namespaceId, userId, uuid, inputProxy)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Update monitor partition fail. err: %s", err))
		return
	}

	data := common.Map[string, interface{}]{}
	data["info"] = info
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (m *monitorController) delPartition(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")

	err := m.monitorService.DelPartition(ginCtx, namespaceId, uuid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Delete external-app fail. err: %s", err))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 请求/转发统计
func (m *monitorController) overviewSummary(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	input := new(monitor_dto.MonSummaryInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Start <= 0 || input.End <= 0 || input.Start >= input.End {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("timestamp is illegal."))
		return
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.UUid)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	wheres := make([]model.MonWhereItem, 0, len(input.Clusters))
	if len(clusterIds) > 0 {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "cluster",
			Operation: "in",
			Values:    clusterIds,
		})
	}

	request, proxy, err := m.monStatistics.CircularMap(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	resMap := make(map[string]interface{})
	resMap["request_summary"] = monitor_dto.CircularDate{
		Total:     request.Total,
		Success:   request.Success,
		Fail:      request.Fail,
		Status4Xx: request.Status4Xx,
		Status5Xx: request.Status5Xx,
	}
	resMap["proxy_summary"] = monitor_dto.CircularDate{
		Total:     proxy.Total,
		Success:   proxy.Success,
		Fail:      proxy.Fail,
		Status4Xx: proxy.Status4Xx,
		Status5Xx: proxy.Status5Xx,
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))

}

// 调用量统计
func (m *monitorController) overviewInvoke(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	input := new(monitor_dto.MonSummaryInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Start <= 0 || input.End <= 0 || input.Start >= input.End {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("timestamp is illegal."))
		return
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.UUid)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	wheres := make([]model.MonWhereItem, 0, len(input.Clusters))
	if len(clusterIds) > 0 {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "cluster",
			Operation: "in",
			Values:    clusterIds,
		})
	}

	values, timeInterval, err := m.monStatistics.Trend(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resMap := make(map[string]interface{})

	resMap["date"] = values.Date
	resMap["request_total"] = values.RequestTotal
	resMap["request_rate"] = values.RequestRate
	resMap["proxy_total"] = values.ProxyTotal
	resMap["proxy_rate"] = values.ProxyRate
	resMap["status_4xx"] = values.Status4XX
	resMap["status_5xx"] = values.Status5XX
	resMap["time_interval"] = timeInterval
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// 报文量统计
func (m *monitorController) overviewMessage(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	input := new(monitor_dto.MonSummaryInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Start <= 0 || input.End <= 0 || input.Start >= input.End {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("timestamp is illegal."))
		return
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.UUid)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	wheres := make([]model.MonWhereItem, 0, len(input.Clusters))
	if len(clusterIds) > 0 {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "cluster",
			Operation: "in",
			Values:    clusterIds,
		})
	}

	data, timeInterval, err := m.monStatistics.MessageTrend(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resMap := make(map[string]interface{})
	resMap["date"] = data.Dates
	resMap["request"] = data.ReqMessage
	resMap["response"] = data.RespMessage
	resMap["time_interval"] = timeInterval
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// API/应用/上游调用量TOP10
func (m *monitorController) overviewTop(ginCtx *gin.Context) {
	input := new(monitor_dto.MonSummaryInput)
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Start <= 0 || input.End <= 0 || input.Start >= input.End {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("timestamp is illegal."))
		return
	}

	//校验data_type
	switch input.DataType {
	case "api", "app", "service":
	default:
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("data_type is illegal."))
		return
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.UUid)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	wheres := make([]model.MonWhereItem, 0, len(input.Clusters))
	if len(clusterIds) > 0 {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "cluster",
			Operation: "in",
			Values:    clusterIds,
		})
	}

	resMap := make(map[string]interface{})
	switch input.DataType {
	case "api":
		//api TOP10
		apiList, err := m.monStatistics.Statistics(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), "api", wheres, 10)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		//补全api信息
		apiUUIDS := make([]string, 0, len(apiList))
		for key := range apiList {
			apiUUIDS = append(apiUUIDS, key)
		}
		apiItems, _ := m.apiService.GetAPIInfoByUUIDS(ginCtx, namespaceId, apiUUIDS)
		apiMaps := common.SliceToMap(apiItems, func(t *api_model.APIInfo) string {
			return t.UUID
		})
		apiResults := make([]*monitor_dto.MonCommonStatisticsOutput, 0, len(apiList))
		for key, val := range apiList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ApiId:         key,
				MonCommonData: newMonCommonData(val),
			}
			if key == "-" {
				monCommonStatisticsOutput.ApiId = "-"
				monCommonStatisticsOutput.ApiName = "无API"
				monCommonStatisticsOutput.IsRed = true
			} else {
				if apiItem, has := apiMaps[key]; has {
					monCommonStatisticsOutput.ApiName = apiItem.Name
					monCommonStatisticsOutput.Path = apiItem.RequestPath
				} else {
					monCommonStatisticsOutput.ApiName = "未知API-" + key
					monCommonStatisticsOutput.IsRed = true
				}
			}

			apiResults = append(apiResults, monCommonStatisticsOutput)
		}
		//排序api
		sort.Slice(apiResults, func(i, j int) bool {
			return apiResults[i].RequestTotal > apiResults[j].RequestTotal
		})

		resMap["api"] = apiResults

	case "app":
		//app TOP10
		appList, err := m.monStatistics.Statistics(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), "app", wheres, 10)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		//补全应用信息
		appUUIDs := make([]string, 0, len(appList))
		for key := range appList {
			appUUIDs = append(appUUIDs, key)
		}
		appItems, _ := m.appService.AppListByUUIDS(ginCtx, namespaceId, appUUIDs)
		appMaps := common.SliceToMap(appItems, func(t *application_model.Application) string {
			return t.IdStr
		})

		appResults := make([]*monitor_dto.MonCommonStatisticsOutput, 0, len(appList))
		for key, val := range appList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				AppId:         key,
				MonCommonData: newMonCommonData(val),
			}
			if key == "-" {
				monCommonStatisticsOutput.AppId = "-"
				monCommonStatisticsOutput.AppName = "无应用"
				monCommonStatisticsOutput.IsRed = true
			} else {
				appInfo, has := appMaps[key]
				if has {
					monCommonStatisticsOutput.AppName = appInfo.Name
				} else {
					monCommonStatisticsOutput.AppName = "未知应用-" + key
					monCommonStatisticsOutput.IsRed = true
				}
			}

			appResults = append(appResults, monCommonStatisticsOutput)
		}
		//排序app
		sort.Slice(appResults, func(i, j int) bool {
			return appResults[i].RequestTotal > appResults[j].RequestTotal
		})
		resMap["app"] = appResults
	case "service":
		//upstream TOP10
		upstreamList, err := m.monStatistics.ProxyStatistics(ginCtx, namespaceId, input.UUid, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), "upstream", wheres, 10)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}

		//补全上游信息
		serviceNames := make([]string, 0, len(upstreamList))
		for key := range upstreamList {
			serviceNames = append(serviceNames, key)
		}
		serviceItems, _ := m.service.GetServiceListByNames(ginCtx, namespaceId, serviceNames)
		serviceMaps := common.SliceToMap(serviceItems, func(t *upstream_model.ServiceListItem) string {
			return t.Name
		})

		serviceResults := make([]*monitor_dto.MonCommonStatisticsOutput, 0, len(upstreamList))
		for key, val := range upstreamList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ServiceName:   key,
				MonCommonData: newMonCommonData(val),
			}
			_, has := serviceMaps[key]
			if !has {
				monCommonStatisticsOutput.IsRed = true
			}
			serviceResults = append(serviceResults, monCommonStatisticsOutput)
		}
		//排序upstream
		sort.Slice(serviceResults, func(i, j int) bool {
			return serviceResults[i].ProxyTotal > serviceResults[j].ProxyTotal
		})

		resMap["service"] = serviceResults
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

func (m *monitorController) getStatistics(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	dataType := ginCtx.Param("dataType")
	detailsType := ginCtx.Param("detailsType")

	resList := make([]*monitor_dto.MonCommonStatisticsOutput, 0)
	resMap := make(map[string]interface{})

	input := &monitor_dto.MonCommonInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if input.EndTime > time.Now().Unix() {
		controller.ErrorJson(ginCtx, http.StatusOK, "查询结束时间不能大于当前时间")
		return
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.PartitionId)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	apiList := make([]*api_model.APIInfo, 0)
	appList := make([]*application_model.Application, 0)
	serviceList := make([]*upstream_model.ServiceListItem, 0)

	callType := "apiOrApp"

	groupBy := ""
	switch dataType {
	case "api":
		groupBy = "api"
		if len(input.ApiIds) > 0 {
			apiList, err = m.apiService.GetAPIInfoByUUIDS(ginCtx, namespaceId, input.ApiIds)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apis error = "+err.Error())
				return
			}
		} else if len(input.Services) > 0 {

			apiInfoList, err := m.apiService.GetAPIListByServiceName(ginCtx, namespaceId, input.Services)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apis error = "+err.Error())
				return
			}
			appIds := common.SliceToSliceIds(apiInfoList, func(t *api_model.APIInfo) string {
				return t.UUID
			})

			apiList, err = m.apiService.GetAPIInfoByUUIDS(ginCtx, namespaceId, appIds)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apis error = "+err.Error())
				return
			}
		} else if detailsType == "" {
			apiList, err = m.apiService.GetAPIInfoAll(ginCtx, namespaceId)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apis error = "+err.Error())
				return
			}
		}

	case "ip":
		groupBy = "ip"
		callType = "ip"
	case "app":
		groupBy = "app"
		if len(input.AppIds) > 0 {
			appList, err = m.appService.AppListByUUIDS(ginCtx, namespaceId, input.AppIds)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apps error = "+err.Error())
				return
			}
		} else if detailsType == "" {
			appList, err = m.appService.AppListAll(ginCtx, namespaceId)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apps error = "+err.Error())
				return
			}
		}

	case "service":
		groupBy = "upstream"
		callType = "proxy"
		if len(input.Services) > 0 {
			serviceList, err = m.service.GetServiceListByNames(ginCtx, namespaceId, input.Services)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query services error = "+err.Error())
				return
			}
		} else {
			serviceList, err = m.service.GetServiceListAll(ginCtx, namespaceId, "")
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query services error = "+err.Error())
				return
			}
		}

	default:
		controller.ErrorJson(ginCtx, http.StatusOK, "参数错误")
		return
	}

	switch detailsType {
	case "api":
		groupBy = "api"
	case "app":
		groupBy = "app"
	case "ip":
		groupBy = "ip"
	case "path", "proxy_path":
		groupBy = "path"
	case "addr":
		groupBy = "addr"
	}

	wheres, err := m.getWheres(ginCtx, namespaceId, input, clusterIds)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	maps := make(map[string]model.MonCommonData)

	if callType == "apiOrApp" {
		maps, err = m.monStatistics.Statistics(ginCtx, namespaceId, input.PartitionId, formatTimeByMinute(input.StartTime), formatTimeByMinute(input.EndTime), groupBy, wheres, 0)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
	} else if callType == "proxy" {
		maps, err = m.monStatistics.ProxyStatistics(ginCtx, namespaceId, input.PartitionId, formatTimeByMinute(input.StartTime), formatTimeByMinute(input.EndTime), groupBy, wheres, 0)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
	}
	//else if callType == "ip" {
	//	maps, err = m.monStatistics.IPStatistics(ginCtx, namespaceId, input.PartitionId, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), groupBy, wheres, 200)
	//	if err != nil {
	//		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
	//		return
	//	}
	//}

	switch detailsType {
	case "api":
		if len(maps) > 0 {
			apiIds := make([]string, 0)
			for key, _ := range maps {
				apiIds = append(apiIds, key)
			}
			apiList, err = m.apiService.GetAPIInfoByUUIDS(ginCtx, namespaceId, apiIds)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query apis error = "+err.Error())
				return
			}

		}

	case "app":
		if len(maps) > 0 {
			appIds := make([]string, 0)
			for key, _ := range maps {
				appIds = append(appIds, key)
			}
			appList, err = m.appService.AppListByUUIDS(ginCtx, namespaceId, appIds)
			if err != nil {
				controller.ErrorJson(ginCtx, http.StatusOK, "query app error = "+err.Error())
				return
			}
		}

	}

	newApiList := make([]*api_model.APIInfo, 0)
	if input.Path != "" {
		for _, info := range apiList {
			if info.RequestPathLabel == input.Path {
				newApiList = append(newApiList, info)
			}
		}
	} else {
		newApiList = apiList
	}

	switch groupBy {
	case "api":
		for _, info := range newApiList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ApiId:   info.UUID,
				ApiName: info.Name,
				//ServiceName: info.ServiceName,
				Path: info.RequestPathLabel,
			}
			if val, ok := maps[info.UUID]; ok {
				monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)

				delete(maps, info.UUID)
			}
			resList = append(resList, monCommonStatisticsOutput)
		}
		for key, val := range maps {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ApiId:   key,
				ApiName: "未知API-" + key,
				IsRed:   true,
			}
			if key == "-" {
				monCommonStatisticsOutput.ApiId = "-"
				monCommonStatisticsOutput.ApiName = "无API"
			}

			monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
			resList = append(resList, monCommonStatisticsOutput)
		}

	case "app":
		for _, info := range appList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				AppId:   info.IdStr,
				AppName: info.Name,
			}
			if val, ok := maps[info.IdStr]; ok {
				monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
				delete(maps, info.IdStr)
			}
			resList = append(resList, monCommonStatisticsOutput)
		}
		for key, val := range maps {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				AppId:   key,
				AppName: "未知应用-" + key,
				IsRed:   true,
			}
			if key == "-" {
				monCommonStatisticsOutput.AppId = "-"
				monCommonStatisticsOutput.AppName = "无应用"
			}

			monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
			resList = append(resList, monCommonStatisticsOutput)
		}
	case "upstream":
		for _, info := range serviceList {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ServiceName: info.Name,
			}
			if val, ok := maps[info.Name]; ok {
				monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
				delete(maps, info.Name)
			}
			resList = append(resList, monCommonStatisticsOutput)
		}
		for _, val := range maps {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{
				ServiceName: "未知服务",
				IsRed:       true,
			}
			monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
			resList = append(resList, monCommonStatisticsOutput)
		}
	case "path", "addr", "ip":
		for key, val := range maps {
			monCommonStatisticsOutput := &monitor_dto.MonCommonStatisticsOutput{}
			if groupBy == "path" {
				monCommonStatisticsOutput.Path = key
				monCommonStatisticsOutput.ProxyPath = key
			} else if groupBy == "addr" {
				monCommonStatisticsOutput.Addr = key
			} else if groupBy == "ip" {
				monCommonStatisticsOutput.Ip = key
			}

			monCommonStatisticsOutput.MonCommonData = newMonCommonData(val)
			resList = append(resList, monCommonStatisticsOutput)
		}

	}

	sort.Slice(resList, func(i, j int) bool {
		return resList[i].RequestTotal > resList[j].RequestTotal
	})
	if groupBy == "upstream" {
		sort.Slice(resList, func(i, j int) bool {
			return resList[i].ProxyTotal > resList[j].ProxyTotal
		})
	}

	resMap["statistics"] = resList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

func (m *monitorController) getStatisticsDetails(ginCtx *gin.Context) {
	dataType := ginCtx.Param("dataType")
	detailsType := ginCtx.Param("detailsType")

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	resMap := make(map[string]interface{})

	input := &monitor_dto.MonCommonInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.EndTime > time.Now().Unix() {
		controller.ErrorJson(ginCtx, http.StatusOK, "查询结束时间不能大于当前时间")
		return
	}

	//校验detailsType
	switch detailsType {
	case "api":
		if input.ApiId == "" {
			controller.ErrorJson(ginCtx, http.StatusOK, "api_id can't be null. ")
			return
		}
	case "app":
		if input.AppId == "" {
			controller.ErrorJson(ginCtx, http.StatusOK, "app_id can't be null. ")
			return
		}
	case "addr":
		if input.Addr == "" {
			controller.ErrorJson(ginCtx, http.StatusOK, "addr can't be null. ")
			return
		}
	}

	if len(input.Clusters) == 0 {
		info, err := m.monitorService.PartitionInfo(ginCtx, namespaceId, input.PartitionId)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		input.Clusters = info.ClusterNames
	}

	clusters, err := m.clusterService.GetByNames(ginCtx, namespaceId, input.Clusters)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	clusterIds := common.SliceToSliceIds(clusters, func(t *cluster_model.Cluster) string {
		return t.UUID
	})

	wheres, err := m.getWheres(ginCtx, namespaceId, input, clusterIds)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	values := &model.MonCallCountInfo{}

	timeInterval := ""
	if dataType == "service" {
		values, timeInterval, err = m.monStatistics.ProxyTrend(ginCtx, namespaceId, input.PartitionId, formatTimeByMinute(input.StartTime), formatTimeByMinute(input.EndTime), wheres)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
	} else {
		values, timeInterval, err = m.monStatistics.Trend(ginCtx, namespaceId, input.PartitionId, formatTimeByMinute(input.StartTime), formatTimeByMinute(input.EndTime), wheres)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
	}

	resValue := &monitor_dto.MonCallCountOutput{
		Date:         values.Date,
		Status5XX:    values.Status5XX,
		Status4XX:    values.Status4XX,
		ProxyRate:    values.ProxyRate,
		ProxyTotal:   values.ProxyTotal,
		RequestRate:  values.RequestRate,
		RequestTotal: values.RequestTotal,
	}

	resMap["tendency"] = resValue
	resMap["time_interval"] = timeInterval

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

func (m *monitorController) getWheres(ginCtx *gin.Context, namespaceId int, input *monitor_dto.MonCommonInput, clusterIds []string) ([]model.MonWhereItem, error) {
	wheres := make([]model.MonWhereItem, 0)

	if input.Ip != "" {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "ip",
			Operation: "=",
			Values:    []string{input.Ip},
		})
	}

	if input.Path != "" {
		apiList, err := m.apiService.GetAPIInfoByPath(ginCtx, namespaceId, input.Path)
		if err != nil {
			return nil, errors.New("query apps error = " + err.Error())
		}
		apiIds := make([]string, 0, len(apiList))
		if len(apiList) > 0 {
			for _, api := range apiList {
				apiIds = append(apiIds, api.UUID)
			}
		} else {
			//若查询的path没有对应的api
			apiIds = append(apiIds, "")
		}

		wheres = append(wheres, model.MonWhereItem{
			Key:       "api",
			Operation: "in",
			Values:    apiIds,
		})
	}

	if input.AppId != "" {
		value := input.AppId
		if input.AppId == "-" {
			value = "-"
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "app",
			Operation: "=",
			Values:    []string{value},
		})
	}

	if input.ProxyPath != "" {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "path",
			Operation: "=",
			Values:    []string{input.ProxyPath},
		})
	}

	if input.Addr != "" {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "addr",
			Operation: "=",
			Values:    []string{input.Addr},
		})
	}

	if input.ServiceName != "" {
		value := input.ServiceName
		if input.ServiceName == "未知服务" {
			value = "-"
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "upstream",
			Operation: "=",
			Values:    []string{value},
		})
	}

	if input.ApiId != "" {
		value := input.ApiId
		if input.ApiId == "-" {
			value = "-"
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "api",
			Operation: "in",
			Values:    []string{value},
		})
	}

	if len(clusterIds) > 0 {
		wheres = append(wheres, model.MonWhereItem{
			Key:       "cluster",
			Operation: "in",
			Values:    clusterIds,
		})
	}

	if len(input.Services) > 0 {
		services := make([]string, 0, len(input.Services))
		for _, v := range input.Services {
			services = append(services, v)
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "upstream",
			Operation: "in",
			Values:    services,
		})
	}

	if len(input.AppIds) > 0 {
		apps := make([]string, 0, len(input.AppIds))
		for _, v := range input.AppIds {
			apps = append(apps, v)
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "app",
			Operation: "in",
			Values:    apps,
		})
	}

	if len(input.ApiIds) > 0 {
		apis := make([]string, 0, len(input.ApiIds))
		for _, v := range input.ApiIds {
			apis = append(apis, v)
		}
		wheres = append(wheres, model.MonWhereItem{
			Key:       "api",
			Operation: "in",
			Values:    apis,
		})
	}

	return wheres, nil
}

func formatTimeByMinute(org int64) time.Time {
	t := time.Unix(org, 0)
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, location)
}

func newMonCommonData(input model.MonCommonData) monitor_dto.MonCommonData {
	return monitor_dto.MonCommonData{
		RequestTotal:   input.RequestTotal,
		RequestSuccess: input.RequestSuccess,
		RequestRate:    input.RequestRate,
		ProxyTotal:     input.ProxyTotal,
		ProxySuccess:   input.ProxySuccess,
		ProxyRate:      input.ProxyRate,
		StatusFail:     input.StatusFail,
		AvgResp:        common.FormatFloat64(input.AvgResp, 2),
		MaxResp:        input.MaxResp,
		MinResp:        input.MinResp,
		AvgTraffic:     common.FormatFloat64(input.AvgTraffic/1024, 2),
		MaxTraffic:     common.FormatFloat64(float64(input.MaxTraffic)/1024, 2),
		MinTraffic:     common.FormatFloat64(float64(input.MinTraffic)/1024, 2),
	}
}
