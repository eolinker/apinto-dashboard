package controller

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	service "github.com/eolinker/apinto-dashboard/modules/api"
	api_model "github.com/eolinker/apinto-dashboard/modules/api/model"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/monitor"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"
	"github.com/eolinker/apinto-dashboard/modules/warn"
	warn_dto "github.com/eolinker/apinto-dashboard/modules/warn/warn-dto"
	warn_model "github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type warnController struct {
	noticeChannelService notice.INoticeChannelService
	warnHistoryService   warn.IWarnHistoryService
	warnStrategyService  warn.IWarnStrategyService
	monitorService       monitor.IMonitorService
	apiService           service.IAPIService
	service              upstream.IService
}

func RegisterWarnRouter(router gin.IRouter) {
	w := &warnController{}
	bean.Autowired(&w.noticeChannelService)
	bean.Autowired(&w.warnStrategyService)
	bean.Autowired(&w.warnHistoryService)
	bean.Autowired(&w.monitorService)
	bean.Autowired(&w.apiService)
	bean.Autowired(&w.service)

	prefix := "/warn"

	//告警历史
	router.GET(prefix+"/history", w.warnHistory)

	//告警策略
	router.GET(prefix+"/strategys", w.strategys)
	router.POST(prefix+"/strategy", w.createStrategy)
	router.PUT(prefix+"/strategy", w.updateStrategy)
	router.GET(prefix+"/strategy", w.strategy)
	router.PATCH(prefix+"/strategy", w.updateStrategyStatus)
	router.DELETE(prefix+"/strategy", w.deleteStrategy)
}

// warnHistory 告警历史
func (w *warnController) warnHistory(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	partitionId := ginCtx.Query("partition_id")
	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, partitionId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	start, _ := strconv.Atoi(ginCtx.Query("start_time"))
	end, _ := strconv.Atoi(ginCtx.Query("end_time"))
	if int64(end) > time.Now().Unix() {
		controller.ErrorJson(ginCtx, http.StatusOK, "查询结束时间不能大于当前时间")
		return
	}
	startTime := time.Unix(int64(start), 0)
	endTime := time.Unix(int64(end), 0)

	pageNum, _ := strconv.Atoi(ginCtx.Query("page_num"))
	pageSize, _ := strconv.Atoi(ginCtx.Query("page_size"))
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	name := ginCtx.Query("strategy_name")

	list, total, err := w.warnHistoryService.QueryList(ginCtx, namespaceId, partitionInfo.Id, pageNum, pageSize, startTime, endTime, name)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	history := make([]*warn_dto.WarnHistory, 0, len(list))
	for _, info := range list {
		target := ""
		switch info.Dimension {
		case "cluster", "partition":
			target = "集群："
		case "service":
			target = "上游："
		case "api":
			target = "API："

		}
		history = append(history, &warn_dto.WarnHistory{
			StrategyTitle: info.StrategyTitle,
			WarnTarget:    target + info.Target,
			WarnContent:   info.Content,
			CreateTime:    common.TimeToStr(info.CreateTime),
			Status:        info.Status,
			ErrMsg:        info.ErrMsg,
		})
	}

	resMap := make(map[string]interface{})

	resMap["history"] = history
	resMap["total"] = total
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// strategys 告警策略列表
func (w *warnController) strategys(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	partitionId := ginCtx.Query("partition_id")
	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, partitionId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	start, _ := strconv.Atoi(ginCtx.Query("start_time"))
	end, _ := strconv.Atoi(ginCtx.Query("end_time"))
	if int64(end) > time.Now().Unix() {
		controller.ErrorJson(ginCtx, http.StatusOK, "查询结束时间不能大于当前时间")
		return
	}
	startTime := time.Unix(int64(start), 0)
	endTime := time.Unix(int64(end), 0)

	status, _ := strconv.Atoi(ginCtx.Query("status"))
	pageNum, _ := strconv.Atoi(ginCtx.Query("page_num"))
	pageSize, _ := strconv.Atoi(ginCtx.Query("page_size"))
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	dimension := ginCtx.Query("warn_dimension")
	dimensions := make([]string, 0)
	if dimension != "" {
		dimensions = strings.Split(dimension, ",")
	}

	query := &warn_model.QueryWarnStrategyParam{
		PartitionId:  partitionInfo.Id,
		StartTime:    startTime,
		EndTime:      endTime,
		StrategyName: ginCtx.Query("strategy_name"),
		Dimension:    dimensions, //告警维度 多选,分割   api/service/cluster/partition
		Status:       status,
		PageNum:      pageNum,
		PageSize:     pageSize,
	}
	listPage, count, err := w.warnStrategyService.WarnStrategyListPage(ginCtx, namespaceId, query)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//获取所有API
	apiList, err := w.apiService.GetAPIInfoAll(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	apiMaps := common.SliceToMap(apiList, func(t *api_model.APIInfo) string {
		return t.UUID
	})

	//获取所有上游服务
	serviceList, err := w.service.GetServiceListAll(ginCtx, namespaceId, "")
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	serviceMaps := common.SliceToMap(serviceList, func(t *upstream_model.ServiceListItem) string {
		return t.Name
	})

	strategyList := make([]*warn_dto.WarnStrategyList, 0, len(listPage))
	for _, strategy := range listPage {

		warnTarget := ""
		targetRule := strategy.WarnStrategyConfig.Target.Rule
		targetValues := strategy.WarnStrategyConfig.Target.Values
		strategyDimension := ""

		switch strategy.Dimension {
		case "api":
			if targetRule == "unlimited" {
				warnTarget = "所有API"
			} else if targetRule == "contain" {
				apiNames := make([]string, 0)
				for _, v := range targetValues {
					if api, ok := apiMaps[v]; ok {
						apiNames = append(apiNames, api.Name)
					}
				}
				sort.Strings(apiNames)
				warnTarget = fmt.Sprintf("API:%s", strings.Join(apiNames, ","))
			} else if targetRule == "not_contain" {

				tempMaps := common.CopyMaps(apiMaps)
				for _, v := range targetValues {
					delete(tempMaps, v)
				}

				apiNames := make([]string, 0)
				for _, info := range tempMaps {
					apiNames = append(apiNames, info.Name)
				}
				sort.Strings(apiNames)
				warnTarget = fmt.Sprintf("API:%s", strings.Join(apiNames, ","))
			}
			strategyDimension = "按API"
		case "service":
			if targetRule == "unlimited" {
				warnTarget = "所有上游服务"
			} else if targetRule == "contain" {
				serviceNames := make([]string, 0)
				for _, v := range targetValues {
					if info, ok := serviceMaps[v]; ok {
						serviceNames = append(serviceNames, info.Name)
					}
				}
				sort.Strings(serviceNames)
				warnTarget = fmt.Sprintf("上游服务:%s", strings.Join(serviceNames, ","))
			} else if targetRule == "not_contain" {
				tempMaps := common.CopyMaps(apiMaps)

				for _, v := range targetValues {
					delete(tempMaps, v)
				}

				serviceNames := make([]string, 0)
				for _, info := range tempMaps {
					serviceNames = append(serviceNames, info.Name)
				}

				sort.Strings(serviceNames)
				warnTarget = fmt.Sprintf("上游服务:%s", strings.Join(serviceNames, ","))
			}
			strategyDimension = "按上游"
		case "cluster":
			warnTarget = fmt.Sprintf("集群:%s", strings.Join(targetValues, ","))
			strategyDimension = "按集群"
		case "partition":
			warnTarget = "当前分区"
			strategyDimension = "按分区"
		}

		warnFrequency := fmt.Sprintf("%d分钟", strategy.Every)
		if strategy.Every == 60 {
			warnFrequency = "1小时"
		}

		strategyList = append(strategyList, &warn_dto.WarnStrategyList{
			Uuid:          strategy.Uuid,
			StrategyTitle: strategy.Title,
			WarnDimension: strategyDimension,
			WarnTarget:    warnTarget,
			WarnRule:      warn_model.QuotaRuleMap[strategy.Quota],
			WarnFrequency: warnFrequency,
			IsEnable:      strategy.IsEnable,
			Operator:      strategy.Operator,
			UpdateTime:    common.TimeToStr(strategy.UpdateTime),
			CreateTime:    common.TimeToStr(strategy.CreateTime),
		})
	}

	resMap := make(map[string]interface{})

	resMap["datas"] = strategyList
	resMap["total"] = count
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// strategy 获取单个告警策略
func (w *warnController) strategy(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uid := ginCtx.Query("uuid")
	if uid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}
	warnStrategy, err := w.warnStrategyService.WarnStrategyByUuid(ginCtx, namespaceId, uid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resMap := make(map[string]interface{})

	resMap["strategy"] = modelWarnStrategyToDto(warnStrategy)
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// updateStrategyStatus 修改告警策略启用状态
func (w *warnController) updateStrategyStatus(ginCtx *gin.Context) {

	strategy := &warn_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if strategy.Uuid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	if err := w.warnStrategyService.UpdateWarnStrategyStatus(ginCtx, strategy.Uuid, strategy.IsEnable); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (w *warnController) deleteStrategy(ginCtx *gin.Context) {

	uid := ginCtx.Query("uuid")
	if uid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	if err := w.warnStrategyService.DeleteWarnStrategy(ginCtx, uid); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// createStrategy 创建告警策略
func (w *warnController) createStrategy(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	strategy := &warn_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := checkWarnStrategyParam(strategy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, strategy.PartitionId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	strategy.Uuid = uuid.New()

	input := dtoWarnStrategyToModel(strategy, partitionInfo.Id)

	if err := w.warnStrategyService.CreateWarnStrategy(ginCtx, namespaceId, userId, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// updateStrategy 修改告警策略
func (w *warnController) updateStrategy(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	strategy := &warn_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := checkWarnStrategyParam(strategy); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, strategy.PartitionId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	input := dtoWarnStrategyToModel(strategy, partitionInfo.Id)

	if err = w.warnStrategyService.UpdateWarnStrategy(ginCtx, namespaceId, userId, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func checkWarnStrategyParam(input *warn_dto.WarnStrategyInput) error {
	if input.Title == "" {
		return errors.New("标题不能为空")
	}
	if input.PartitionId == "" {
		return errors.New("partition_id is null")
	}
	if input.Dimension != "api" && input.Dimension != "service" && input.Dimension != "cluster" && input.Dimension != "partition" {
		return errors.New("请选择告警维度")
	}
	if input.Quota == "" {
		return errors.New("请选择告警指标")
	}
	if _, ok := warn_model.QuotaRuleMap[warn_model.QuotaType(input.Quota)]; !ok {
		return errors.New("告警指标参数错误")
	}

	if input.Every != 1 && input.Every != 3 && input.Every != 5 && input.Every != 10 && input.Every != 30 && input.Every != 60 {
		return errors.New("统计时间粒度参数错误")
	}
	//switch input.Target.Rule {
	//case "unlimited":
	//	switch input.Dimension {
	//	case "api":
	//		if input.Target.Values {
	//
	//		}
	//	}
	//case "contain":
	//case "not_contain":
	//default:
	//	return errors.New("请选择告警目标")
	//}

	return nil
}

func modelWarnStrategyToDto(strategy *warn_model.WarnStrategy) *warn_dto.WarnStrategy {

	rules := make([]*warn_dto.WarnStrategyRule, 0, len(strategy.WarnStrategyConfig.Rule))
	for _, rule := range strategy.WarnStrategyConfig.Rule {

		condition := make([]*warn_dto.WarnStrategyRuleCondition, 0, len(rule.Condition))

		for _, ruleCondition := range rule.Condition {

			condition = append(condition, &warn_dto.WarnStrategyRuleCondition{
				Compare: ruleCondition.Compare,
				Unit:    ruleCondition.Unit,
				Value:   ruleCondition.Value,
			})
		}

		rules = append(rules, &warn_dto.WarnStrategyRule{
			ChannelUuids: rule.ChannelUuids,
			Condition:    condition,
		})
	}
	return &warn_dto.WarnStrategy{
		Uuid:      strategy.Uuid,
		Title:     strategy.Title,
		Desc:      strategy.Desc,
		IsEnable:  strategy.IsEnable,
		Dimension: strategy.Dimension,
		Target: &warn_dto.WarnStrategyTarget{
			Rule:   strategy.WarnStrategyConfig.Target.Rule,
			Values: strategy.WarnStrategyConfig.Target.Values,
		},
		Quota:      string(strategy.Quota),
		Every:      strategy.Every,
		Rule:       rules,
		Continuity: strategy.WarnStrategyConfig.Continuity,
		HourMax:    strategy.WarnStrategyConfig.HourMax,
		Users:      strategy.WarnStrategyConfig.Users,
	}
}
func dtoWarnStrategyToModel(strategy *warn_dto.WarnStrategyInput, partitionId int) *warn_model.WarnStrategy {
	rule := make([]*warn_model.WarnStrategyConfigRule, 0, len(strategy.Rule))
	for _, strategyRule := range strategy.Rule {
		condition := make([]*warn_model.WarnStrategyConfigRuleCondition, 0)
		for _, ruleCondition := range strategyRule.Condition {
			condition = append(condition, &warn_model.WarnStrategyConfigRuleCondition{
				Compare: ruleCondition.Compare,
				Unit:    ruleCondition.Unit,
				Value:   ruleCondition.Value,
			})
		}
		rule = append(rule, &warn_model.WarnStrategyConfigRule{
			ChannelUuids: strategyRule.ChannelUuids,
			Condition:    condition,
		})
	}

	return &warn_model.WarnStrategy{
		PartitionId: partitionId,
		Uuid:        strategy.Uuid,
		Title:       strategy.Title,
		Desc:        strategy.Desc,
		IsEnable:    strategy.IsEnable,
		Dimension:   strategy.Dimension,
		Quota:       warn_model.QuotaType(strategy.Quota),
		Every:       strategy.Every,
		WarnStrategyConfig: &warn_model.WarnStrategyConfig{
			Target: warn_model.WarnStrategyConfigTarget{
				Rule:   strategy.Target.Rule,
				Values: strategy.Target.Values,
			},
			Rule:       rule,
			Continuity: strategy.Continuity,
			HourMax:    strategy.HourMax,
			Users:      strategy.Users,
		},
		Operator:   "",
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
}
