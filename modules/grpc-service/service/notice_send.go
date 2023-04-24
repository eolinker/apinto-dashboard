package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	notice_model "github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	"github.com/eolinker/apinto-dashboard/modules/user"
	warn_model "github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strings"
	"sync/atomic"
)

var _ grpc_service.NoticeSendServer = (*noticeSendService)(nil)

type noticeSendService struct {
	userService          user.IUserInfoService
	noticeChannelService notice.INoticeChannelService
	noticeChannelDriver  notice.INoticeChannelDriverManager
}

func NewNoticeSendService() grpc_service.NoticeSendServer {
	n := &noticeSendService{}
	bean.Autowired(&n.userService)
	bean.Autowired(&n.noticeChannelService)
	bean.Autowired(&n.noticeChannelDriver)
	return n
}

func (n *noticeSendService) Send(ctx context.Context, req *grpc_service.NoticeSendReq) (*grpc_service.NoticeSendResp, error) {
	//获取用户信息
	userIds := make([]int, 0, len(req.UserIds))
	for _, id := range req.UserIds {
		userIds = append(userIds, int(id))
	}
	userInfos, err := n.userService.GetUserInfoMaps(ctx, userIds...)
	if err != nil {
		log.Errorf("warn-notice send fail. userService.userInfos error:%s", err.Error())
		return nil, fmt.Errorf("获取用户信息失败 error:%s", err)
	}

	channelList, err := n.noticeChannelService.NoticeChannelList(ctx, int(req.NamespaceId), 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("warn-notice send fail. Get NoticeChannelList fail. err=%s", err.Error())
		return nil, fmt.Errorf("获取通知渠道列表失败 error:%s", err)
	}

	noticeChannelMap := common.SliceToMap(channelList, func(t *notice_model.NoticeChannel) string {
		return t.Name
	})

	//获取用户的邮箱和通知渠道ID
	userEmailStr := make([]string, 0)
	noticeUserId := make([]string, 0)
	for _, userId := range userIds {
		if u, ok := userInfos[userId]; ok {
			if len(strings.TrimSpace(u.Email)) > 0 {
				userEmailStr = append(userEmailStr, u.Email)
			}

			if len(strings.TrimSpace(u.NoticeUserId)) > 0 {
				noticeUserId = append(noticeUserId, u.NoticeUserId)
			}
		}
	}

	//发送失败的次数和需要发送的次数做对比
	var sendFail = new(int64)
	noticeErrGroup, _ := errgroup.WithContext(ctx)
	sendMsgErrors := make([]*warn_model.SendMsgError, 0)

	for channelUuid, noticeMsg := range req.Notices {
		//利用协程快速发送通知消息
		noticeErrGroup.Go(func() error {
			noticeChannelDriver := n.noticeChannelDriver.GetDriver(channelUuid)
			if noticeChannelDriver == nil {
				log.Errorf("获取不到通知渠道 渠道uuid：%s", channelUuid)
				return errors.New("渠道通知获取失败")
			}
			sendMsgErrorUuid := uuid.New()
			sendMsgError := &warn_model.SendMsgError{
				UUID:              sendMsgErrorUuid,
				NoticeChannelUUID: channelUuid,
			}

			if channel, ok := noticeChannelMap[channelUuid]; ok {
				sends := make([]string, 0)
				if channel.Type == 2 {
					if len(userEmailStr) == 0 {
						atomic.AddInt64(sendFail, 1)
						sendMsgError.Msg = "收件人邮箱为空"
						sendMsgErrors = append(sendMsgErrors, sendMsgError)
						return errors.New(sendMsgError.Msg)
					}
					//获取邮箱msg
					sends = userEmailStr
				} else {
					if len(noticeUserId) == 0 {
						atomic.AddInt64(sendFail, 1)
						sendMsgError.Msg = "通知用户ID为空"
						sendMsgErrors = append(sendMsgErrors, sendMsgError)
						return errors.New(sendMsgError.Msg)
					}
					sends = noticeUserId
				}
				if err = noticeChannelDriver.SendTo(sends, noticeMsg.Title, noticeMsg.Msg); err != nil {
					sendMsgError.Msg = err.Error()
					sendMsgErrors = append(sendMsgErrors, sendMsgError)
					log.Errorf("告警消息发送失败 sendMsgErrorUuid=%s channelUuid=%s users=%v err=%s", sendMsgErrorUuid, channelUuid, sends, err.Error())
					atomic.AddInt64(sendFail, 1)
					return err
				}
			}

			return nil
		})

	}

	sendStatus := 1
	//发生错误 发送失败
	if err = noticeErrGroup.Wait(); err != nil {
		//部分发送成功
		sendStatus = 3
		//全部发送失败
		if *sendFail == int64(len(req.Notices)) {
			sendStatus = 2
		}
	}
	errMsg, _ := json.Marshal(sendMsgErrors)

	return &grpc_service.NoticeSendResp{
		SendStatus: int32(sendStatus),
		ErrMsg:     string(errMsg),
	}, nil
}
