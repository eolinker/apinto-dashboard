package controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	"github.com/eolinker/apinto-dashboard/modules/notice/dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type noticeController struct {
	noticeChannelService notice.INoticeChannelService
}

func newNoticeController() *noticeController {
	w := &noticeController{}
	bean.Autowired(&w.noticeChannelService)
	return w
}

// channels 可选渠道列表
func (w *noticeController) channels(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	list, err := w.noticeChannelService.NoticeChannelList(ginCtx, namespaceId, 0)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	channel := make([]*dto.NoticeChannel, 0, len(list))
	for _, noticeChannel := range list {
		channel = append(channel, &dto.NoticeChannel{
			Uuid:  noticeChannel.Name,
			Title: noticeChannel.Title,
			Type:  noticeChannel.Type,
		})
	}

	resMap := make(map[string]interface{})

	resMap["channels"] = channel
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(resMap))
}
