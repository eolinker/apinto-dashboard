package controller

import (
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/email/dto"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	notice_model "github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
)

type emailController struct {
	noticeChannelService notice.INoticeChannelService
}

func newEmailController() *emailController {
	w := &emailController{}
	bean.Autowired(&w.noticeChannelService)
	return w
}

// getEmail 获取通知邮箱
func (w *emailController) getEmail(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 2)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resMap := make(map[string]interface{})
	if len(list) > 0 {

		email := &notice_model.NoticeChannelEmail{}

		if err = json.Unmarshal([]byte(list[0].Config), email); err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}

		emailInfo := &dto.EmailOutput{
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

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}

// createEmail 创建通知邮箱
func (w *emailController) createEmail(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	emailInput := new(dto.EmailInput)

	if err := ginCtx.BindJSON(emailInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//参数校验
	if err := checkEmailParam(emailInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
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
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// updateEmail 修改通知邮箱
func (w *emailController) updateEmail(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)
	emailInput := new(dto.EmailInput)

	if err := ginCtx.BindJSON(emailInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if emailInput.Uuid == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "uuid is null")
		return
	}

	//参数校验
	if err := checkEmailParam(emailInput); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
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
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func checkEmailParam(input *dto.EmailInput) error {
	if input.SmtpUrl == "" {
		return errors.New("smtp_url is null")
	}
	if input.SmtpPort == 0 {
		return errors.New("smtp_port is 0")
	}
	return nil
}
