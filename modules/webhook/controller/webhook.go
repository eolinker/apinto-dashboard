package controller

import (
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	notice_model "github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	"github.com/eolinker/apinto-dashboard/modules/webhook/dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
)

type webhookController struct {
	noticeChannelService notice.INoticeChannelService
}

func newWebhookController() *webhookController {
	w := &webhookController{}
	bean.Autowired(&w.noticeChannelService)
	return w
}

// delWebhook 删除webhook
func (w *webhookController) delWebhook(ginCtx *gin.Context) {
	uid := ginCtx.Query("uuid")
	if uid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	if err := w.noticeChannelService.DeleteNoticeChannel(ginCtx, namespaceId, userId, uid); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// webhook 获取单个webhook信息
func (w *webhookController) webhook(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uid := ginCtx.Query("uuid")
	if uid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	channel, err := w.noticeChannelService.NoticeChannelByName(ginCtx, namespaceId, uid)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	webhook := &notice_model.NoticeChannelWebhook{}

	if err = json.Unmarshal([]byte(channel.Config), webhook); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	webhookOutPut := &dto.WebhookOutput{
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
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// webhook 获取webhook列表
func (w *webhookController) webhooks(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 1)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	webhooks := make([]*dto.WebhooksOutput, 0, len(list))
	for _, channel := range list {

		webhook := &notice_model.NoticeChannelWebhook{}

		if err = json.Unmarshal([]byte(channel.Config), webhook); err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		webhooks = append(webhooks, &dto.WebhooksOutput{
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

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// createWebhook 新建webhook
func (w *webhookController) createWebhook(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	webhookInput := new(dto.WebhookInput)

	if err := ginCtx.BindJSON(webhookInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//参数校验
	if err := checkWebhookParam(webhookInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
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
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// updateWebhook 修改webhook
func (w *webhookController) updateWebhook(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	webhookInput := new(dto.WebhookInput)

	if err := ginCtx.BindJSON(webhookInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if webhookInput.Uuid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	//参数校验
	if err := checkWebhookParam(webhookInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
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
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func checkWebhookParam(webhookInput *dto.WebhookInput) error {
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
