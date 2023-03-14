package monitor_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/dto/monitor-dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model/monitor-model"
	"github.com/eolinker/apinto-dashboard/model/notice-model"
	service2 "github.com/eolinker/apinto-dashboard/modules/api"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"
	"github.com/eolinker/apinto-dashboard/service/monitor-service"
	"github.com/eolinker/apinto-dashboard/service/notice-service"
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
	noticeChannelService notice_service.INoticeChannelService
	warnHistoryService   monitor_service.IWarnHistoryService
	warnStrategyService  monitor_service.IWarnStrategyService
	monitorService       monitor_service.IMonitorService
	apiService           service2.IAPIService
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

	//webhook操作
	router.DELETE(prefix+"/webhook", controller.GenAccessHandler(access.NoticeWebhookView, access.NoticeWebhookEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindNoticeWebhook), w.delWebhook)
	router.POST(prefix+"/webhook", controller.GenAccessHandler(access.NoticeWebhookView, access.NoticeWebhookEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindNoticeWebhook), w.createWebhook)
	router.PUT(prefix+"/webhook", controller.GenAccessHandler(access.NoticeWebhookView, access.NoticeWebhookEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindNoticeWebhook), w.updateWebhook)
	router.GET(prefix+"/webhook", controller.GenAccessHandler(access.NoticeWebhookView, access.NoticeWebhookEdit), w.webhook)
	router.GET(prefix+"/webhooks", controller.GenAccessHandler(access.NoticeWebhookView, access.NoticeWebhookEdit), w.webhooks)

	//邮箱操作
	router.POST(prefix+"/email", controller.GenAccessHandler(access.NoticeEmailView, access.NoticeEmailEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindNoticeEmail), w.createEmail)
	router.PUT(prefix+"/email", controller.GenAccessHandler(access.NoticeEmailView, access.NoticeEmailEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindNoticeWebhook), w.updateEmail)
	router.GET(prefix+"/email", controller.GenAccessHandler(access.NoticeEmailView, access.NoticeEmailEdit), w.getEmail)

	//告警历史
	router.GET(prefix+"/history", controller.GenAccessHandler(access.MonPartitionView), w.warnHistory)
	//可选渠道列表
	router.GET(prefix+"/channels", w.channels)

	//告警策略
	router.GET(prefix+"/strategys", controller.GenAccessHandler(access.MonPartitionView), w.strategys)
	router.POST(prefix+"/strategy", controller.GenAccessHandler(access.MonPartitionView, access.MonPartitionEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindWarnStrategy), w.createStrategy)
	router.PUT(prefix+"/strategy", controller.GenAccessHandler(access.MonPartitionView, access.MonPartitionEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindWarnStrategy), w.updateStrategy)
	router.GET(prefix+"/strategy", w.strategy)
	router.PATCH(prefix+"/strategy", controller.GenAccessHandler(access.MonPartitionView, access.MonPartitionEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindWarnStrategy), w.updateStrategyStatus)
	router.DELETE(prefix+"/strategy", controller.GenAccessHandler(access.MonPartitionView, access.MonPartitionEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindWarnStrategy), w.deleteStrategy)
}

// delWebhook 删除webhook
func (w *warnController) delWebhook(ginCtx *gin.Context) {
	uid := ginCtx.Query("uuid")
	if uid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	if err := w.noticeChannelService.DeleteNoticeChannel(ginCtx, namespaceId, userId, uid); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// webhook 获取单个webhook信息
func (w *warnController) webhook(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	uid := ginCtx.Query("uuid")
	if uid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	channel, err := w.noticeChannelService.NoticeChannelByName(ginCtx, namespaceId, uid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	webhook := &notice_model.NoticeChannelWebhook{}

	if err = json.Unmarshal([]byte(channel.Config), webhook); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	webhookOutPut := &monitor_dto.WebhookOutput{
		Uuid:          channel.Name,
		Title:         channel.Title,
		Desc:          webhook.Desc,
		Url:           webhook.Url,
		Method:        webhook.Method,
		ContentType:   webhook.ContentType,
		NoticeType:    webhook.NoticeType,
		UserSeparator: webhook.UserSeparator,
		Header:        webhook.Header,
		Template:      webhook.Template,
	}

	resMap := make(map[string]interface{})

	resMap["webhook"] = webhookOutPut
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// webhook 获取webhook列表
func (w *warnController) webhooks(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 1)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	webhooks := make([]*monitor_dto.WebhooksOutput, 0, len(list))
	for _, channel := range list {

		webhook := &notice_model.NoticeChannelWebhook{}

		if err = json.Unmarshal([]byte(channel.Config), webhook); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}
		webhooks = append(webhooks, &monitor_dto.WebhooksOutput{
			Uuid:        channel.Name,
			Title:       channel.Title,
			Url:         webhook.Url,
			Method:      webhook.Method,
			ContentType: webhook.ContentType,
			Operator:    channel.Operator,
			UpdateTime:  common.TimeToStr(channel.UpdateTime),
			CreateTime:  common.TimeToStr(channel.CreateTime),
			IsDelete:    channel.IsDelete,
		})
	}
	resMap := make(map[string]interface{})

	resMap["webhooks"] = webhooks

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// createWebhook 新建webhook
func (w *warnController) createWebhook(ginCtx *gin.Context) {

	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	webhookInput := new(monitor_dto.WebhookInput)

	if err := ginCtx.BindJSON(webhookInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	//参数校验
	if err := checkWebhookParam(webhookInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	config := &notice_model.NoticeChannelWebhook{
		Desc:          webhookInput.Desc,
		Url:           webhookInput.Url,
		Method:        webhookInput.Method,
		ContentType:   webhookInput.ContentType,
		NoticeType:    webhookInput.NoticeType,
		UserSeparator: webhookInput.UserSeparator,
		Header:        webhookInput.Header,
		Template:      webhookInput.Template,
	}

	bytes, _ := json.Marshal(config)

	channel := &notice_model.NoticeChannel{
		Name:   uuid.New(),
		Title:  webhookInput.Title,
		Type:   1,
		Config: string(bytes),
	}

	if err := w.noticeChannelService.CreateNoticeChannel(ginCtx, namespaceId, userId, channel); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// updateWebhook 修改webhook
func (w *warnController) updateWebhook(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	webhookInput := new(monitor_dto.WebhookInput)

	if err := ginCtx.BindJSON(webhookInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if webhookInput.Uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	//参数校验
	if err := checkWebhookParam(webhookInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	config := &notice_model.NoticeChannelWebhook{
		Desc:          webhookInput.Desc,
		Url:           webhookInput.Url,
		Method:        webhookInput.Method,
		ContentType:   webhookInput.ContentType,
		NoticeType:    webhookInput.NoticeType,
		UserSeparator: webhookInput.UserSeparator,
		Header:        webhookInput.Header,
		Template:      webhookInput.Template,
	}

	bytes, _ := json.Marshal(config)

	channel := &notice_model.NoticeChannel{
		Name:   webhookInput.Uuid,
		Title:  webhookInput.Title,
		Type:   1,
		Config: string(bytes),
	}

	if err := w.noticeChannelService.UpdateNoticeChannel(ginCtx, namespaceId, userId, channel); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func checkWebhookParam(webhookInput *monitor_dto.WebhookInput) error {
	if webhookInput.Url == "" {
		return errors.New("url is null")
	}

	if webhookInput.Title == "" {
		return errors.New("title is null")
	}

	method := webhookInput.Method
	if method != http.MethodPost && method != http.MethodGet {
		return errors.New("method param fail " + method)
	}

	noticeType := webhookInput.NoticeType
	if noticeType != "single" && noticeType != "many" {
		return errors.New("notice_type param fail " + noticeType)
	}

	contentType := webhookInput.ContentType
	if contentType != "JSON" && contentType != "form-data" {
		return errors.New("content_type param fail " + contentType)
	}
	return nil
}

// getEmail 获取通知邮箱
func (w *warnController) getEmail(ginCtx *gin.Context) {

	namespaceId := controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 2)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resMap := make(map[string]interface{})
	if len(list) > 0 {

		email := &notice_model.NoticeChannelEmail{}

		if err = json.Unmarshal([]byte(list[0].Config), email); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		emailInfo := &monitor_dto.EmailOutput{
			Uuid:     list[0].Name,
			SmtpUrl:  email.SmtpUrl,
			SmtpPort: email.SmtpPort,
			Protocol: email.Protocol,
			Email:    list[0].Title,
			Account:  email.Account,
			Password: email.Password,
		}

		resMap["email_info"] = emailInfo
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// createEmail 创建通知邮箱
func (w *warnController) createEmail(ginCtx *gin.Context) {

	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	emailInput := new(monitor_dto.EmailInput)

	if err := ginCtx.BindJSON(emailInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	//参数校验
	if err := checkEmailParam(emailInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	config := &notice_model.NoticeChannelEmail{
		SmtpUrl:  emailInput.SmtpUrl,
		Email:    emailInput.Email,
		SmtpPort: emailInput.SmtpPort,
		Protocol: emailInput.Protocol,
		Account:  emailInput.Account,
		Password: emailInput.Password,
	}

	bytes, _ := json.Marshal(config)

	channel := &notice_model.NoticeChannel{
		Name:   uuid.New(),
		Title:  emailInput.Email,
		Type:   2,
		Config: string(bytes),
	}

	if err := w.noticeChannelService.CreateNoticeChannel(ginCtx, namespaceId, userId, channel); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// updateEmail 修改通知邮箱
func (w *warnController) updateEmail(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	emailInput := new(monitor_dto.EmailInput)

	if err := ginCtx.BindJSON(emailInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	if emailInput.Uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	//参数校验
	if err := checkEmailParam(emailInput); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	config := &notice_model.NoticeChannelEmail{
		SmtpUrl:  emailInput.SmtpUrl,
		SmtpPort: emailInput.SmtpPort,
		Protocol: emailInput.Protocol,
		Account:  emailInput.Account,
		Email:    emailInput.Email,
		Password: emailInput.Password,
	}

	bytes, _ := json.Marshal(config)

	channel := &notice_model.NoticeChannel{
		Name:   emailInput.Uuid,
		Title:  emailInput.Email,
		Type:   2,
		Config: string(bytes),
	}

	if err := w.noticeChannelService.UpdateNoticeChannel(ginCtx, namespaceId, userId, channel); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func checkEmailParam(input *monitor_dto.EmailInput) error {
	if input.SmtpUrl == "" {
		return errors.New("smtp_url is null")
	}
	if input.SmtpPort == 0 {
		return errors.New("smtp_port is 0")
	}
	return nil
}

// warnHistory 告警历史
func (w *warnController) warnHistory(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)

	partitionId := ginCtx.Query("partition_id")
	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, partitionId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	start, _ := strconv.Atoi(ginCtx.Query("start_time"))
	end, _ := strconv.Atoi(ginCtx.Query("end_time"))
	if int64(end) > time.Now().Unix() {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("查询结束时间不能大于当前时间"))
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
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	history := make([]*monitor_dto.WarnHistory, 0, len(list))
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
		history = append(history, &monitor_dto.WarnHistory{
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
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// channels 可选渠道列表
func (w *warnController) channels(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 0)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	channel := make([]*monitor_dto.NoticeChannel, 0, len(list))
	for _, noticeChannel := range list {
		channel = append(channel, &monitor_dto.NoticeChannel{
			Uuid:  noticeChannel.Name,
			Title: noticeChannel.Title,
			Type:  noticeChannel.Type,
		})
	}

	resMap := make(map[string]interface{})

	resMap["channels"] = channel
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// strategys 告警策略列表
func (w *warnController) strategys(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	partitionId := ginCtx.Query("partition_id")
	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, partitionId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	start, _ := strconv.Atoi(ginCtx.Query("start_time"))
	end, _ := strconv.Atoi(ginCtx.Query("end_time"))
	if int64(end) > time.Now().Unix() {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("查询结束时间不能大于当前时间"))
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

	query := &monitor_model.QueryWarnStrategyParam{
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
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	//获取所有API
	apiList, err := w.apiService.GetAPIInfoAll(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	apiMaps := common.SliceToMap(apiList, func(t *apimodel.APIInfo) string {
		return t.UUID
	})

	//获取所有上游服务
	serviceList, err := w.service.GetServiceListAll(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	serviceMaps := common.SliceToMap(serviceList, func(t *upstream_model.ServiceListItem) string {
		return t.Name
	})

	strategyList := make([]*monitor_dto.WarnStrategyList, 0, len(listPage))
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

		strategyList = append(strategyList, &monitor_dto.WarnStrategyList{
			Uuid:          strategy.Uuid,
			StrategyTitle: strategy.Title,
			WarnDimension: strategyDimension,
			WarnTarget:    warnTarget,
			WarnRule:      monitor_model.QuotaRuleMap[strategy.Quota],
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
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// strategy 获取单个告警策略
func (w *warnController) strategy(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	uid := ginCtx.Query("uuid")
	if uid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}
	warnStrategy, err := w.warnStrategyService.WarnStrategyByUuid(ginCtx, namespaceId, uid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resMap := make(map[string]interface{})

	resMap["strategy"] = modelWarnStrategyToDto(warnStrategy)
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(resMap))
}

// updateStrategyStatus 修改告警策略启用状态
func (w *warnController) updateStrategyStatus(ginCtx *gin.Context) {

	strategy := &monitor_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if strategy.Uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	if err := w.warnStrategyService.UpdateWarnStrategyStatus(ginCtx, strategy.Uuid, strategy.IsEnable); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (w *warnController) deleteStrategy(ginCtx *gin.Context) {

	uid := ginCtx.Query("uuid")
	if uid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("uuid is null"))
		return
	}

	if err := w.warnStrategyService.DeleteWarnStrategy(ginCtx, uid); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// createStrategy 创建告警策略
func (w *warnController) createStrategy(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	strategy := &monitor_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := checkWarnStrategyParam(strategy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, strategy.PartitionId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	strategy.Uuid = uuid.New()

	input := dtoWarnStrategyToModel(strategy, partitionInfo.Id)

	if err := w.warnStrategyService.CreateWarnStrategy(ginCtx, namespaceId, userId, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// updateStrategy 修改告警策略
func (w *warnController) updateStrategy(ginCtx *gin.Context) {

	namespaceId := controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	strategy := &monitor_dto.WarnStrategyInput{}

	if err := ginCtx.BindJSON(strategy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := checkWarnStrategyParam(strategy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	partitionInfo, err := w.monitorService.PartitionInfo(ginCtx, namespaceId, strategy.PartitionId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	input := dtoWarnStrategyToModel(strategy, partitionInfo.Id)

	if err = w.warnStrategyService.UpdateWarnStrategy(ginCtx, namespaceId, userId, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func checkWarnStrategyParam(input *monitor_dto.WarnStrategyInput) error {
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
	if _, ok := monitor_model.QuotaRuleMap[monitor_model.QuotaType(input.Quota)]; !ok {
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

func modelWarnStrategyToDto(strategy *monitor_model.WarnStrategy) *monitor_dto.WarnStrategy {

	rules := make([]*monitor_dto.WarnStrategyRule, 0, len(strategy.WarnStrategyConfig.Rule))
	for _, rule := range strategy.WarnStrategyConfig.Rule {

		condition := make([]*monitor_dto.WarnStrategyRuleCondition, 0, len(rule.Condition))

		for _, ruleCondition := range rule.Condition {

			condition = append(condition, &monitor_dto.WarnStrategyRuleCondition{
				Compare: ruleCondition.Compare,
				Unit:    ruleCondition.Unit,
				Value:   ruleCondition.Value,
			})
		}

		rules = append(rules, &monitor_dto.WarnStrategyRule{
			ChannelUuids: rule.ChannelUuids,
			Condition:    condition,
		})
	}
	return &monitor_dto.WarnStrategy{
		Uuid:      strategy.Uuid,
		Title:     strategy.Title,
		Desc:      strategy.Desc,
		IsEnable:  strategy.IsEnable,
		Dimension: strategy.Dimension,
		Target: &monitor_dto.WarnStrategyTarget{
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
func dtoWarnStrategyToModel(strategy *monitor_dto.WarnStrategyInput, partitionId int) *monitor_model.WarnStrategy {
	rule := make([]*monitor_model.WarnStrategyConfigRule, 0, len(strategy.Rule))
	for _, strategyRule := range strategy.Rule {
		condition := make([]*monitor_model.WarnStrategyConfigRuleCondition, 0)
		for _, ruleCondition := range strategyRule.Condition {
			condition = append(condition, &monitor_model.WarnStrategyConfigRuleCondition{
				Compare: ruleCondition.Compare,
				Unit:    ruleCondition.Unit,
				Value:   ruleCondition.Value,
			})
		}
		rule = append(rule, &monitor_model.WarnStrategyConfigRule{
			ChannelUuids: strategyRule.ChannelUuids,
			Condition:    condition,
		})
	}

	return &monitor_model.WarnStrategy{
		PartitionId: partitionId,
		Uuid:        strategy.Uuid,
		Title:       strategy.Title,
		Desc:        strategy.Desc,
		IsEnable:    strategy.IsEnable,
		Dimension:   strategy.Dimension,
		Quota:       monitor_model.QuotaType(strategy.Quota),
		Every:       strategy.Every,
		WarnStrategyConfig: &monitor_model.WarnStrategyConfig{
			Target: monitor_model.WarnStrategyConfigTarget{
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
