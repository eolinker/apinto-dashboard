package notice_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/notice"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-entry"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
	notice_store2 "github.com/eolinker/apinto-dashboard/modules/notice/notice-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

type noticeChannelService struct {
	noticeChannelStore   notice_store2.INoticeChannelStore
	noticeChannelStat    notice_store2.INoticeChannelStatStore
	noticeChannelVersion notice_store2.INoticeChannelVersionStore
	noticeChannelDriver  notice.INoticeChannelDriverManager
	quoteStore           quote_store.IQuoteStore
	userService          user.IUserInfoService
}

func newNoticeChannelService() notice.INoticeChannelService {
	n := &noticeChannelService{}
	bean.Autowired(&n.noticeChannelStore)
	bean.Autowired(&n.quoteStore)
	bean.Autowired(&n.noticeChannelStat)
	bean.Autowired(&n.noticeChannelVersion)
	bean.Autowired(&n.userService)
	bean.Autowired(&n.noticeChannelDriver)
	return n
}

func (n *noticeChannelService) InitChannelDriver() error {
	ctx := context.Background()
	noticeChannels, err := n.noticeChannelStore.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, channel := range noticeChannels {
		version, err := n.latestNoticeChannelVersion(ctx, channel.Id)
		if err != nil {
			return err
		}
		var noticeChannelDriver notice.IDriverNoticeChannel
		if channel.Type == 2 {
			email := new(notice_model.NoticeChannelEmail)
			if err = json.Unmarshal([]byte(version.Config), email); err != nil {
				return err
			}
			noticeChannelDriver = common.NewSmtp(email.SmtpUrl, email.SmtpPort, email.Protocol, email.Account, email.Password, email.Email)
		} else {
			webhook := new(notice_model.NoticeChannelWebhook)
			if err := json.Unmarshal([]byte(version.Config), webhook); err != nil {
				return err
			}
			noticeChannelDriver = common.NewWebhook(webhook.Url, webhook.Method, webhook.ContentType, webhook.NoticeType, webhook.UserSeparator, webhook.Header, webhook.Template)
		}
		n.noticeChannelDriver.RegisterDriver(channel.Name, noticeChannelDriver)
	}
	return nil
}

func (n *noticeChannelService) CreateNoticeChannel(ctx context.Context, namespaceId, userID int, channel *notice_model.NoticeChannel) error {

	t := time.Now()
	noticeChannel := &notice_entry.NoticeChannel{
		NamespaceID: namespaceId,
		Name:        channel.Name,
		Title:       channel.Title,
		Type:        channel.Type,
		Operator:    userID,
		CreateTime:  t,
		UpdateTime:  t,
	}

	var noticeChannelDriver notice.IDriverNoticeChannel
	//邮箱
	if channel.Type == 2 {
		channels, _ := n.noticeChannelStore.GetByType(ctx, namespaceId, channel.Type)
		if len(channels) > 0 {
			return errors.New("目前仅支持一个邮箱")
		}

		email := new(notice_model.NoticeChannelEmail)
		if err := json.Unmarshal([]byte(channel.Config), email); err != nil {
			return err
		}
		noticeChannelDriver = common.NewSmtp(email.SmtpUrl, email.SmtpPort, email.Protocol, email.Account, email.Password, email.Email)
	} else {
		webhook := new(notice_model.NoticeChannelWebhook)
		if err := json.Unmarshal([]byte(channel.Config), webhook); err != nil {
			return err
		}
		noticeChannelDriver = common.NewWebhook(webhook.Url, webhook.Method, webhook.ContentType, webhook.NoticeType, webhook.UserSeparator, webhook.Header, webhook.Template)
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: channel.Name,
		Name: channel.Title,
	})

	return n.noticeChannelStore.Transaction(ctx, func(txCtx context.Context) error {

		if err := n.noticeChannelStore.Insert(txCtx, noticeChannel); err != nil {
			return err
		}

		noticeChannelVersion := &notice_entry.NoticeChannelVersion{
			NoticeChannelID: noticeChannel.Id,
			NamespaceID:     namespaceId,
			NoticeChannelVersionConfig: notice_entry.NoticeChannelVersionConfig{
				Config: channel.Config,
			},
			Operator:   userID,
			CreateTime: t,
		}

		if err := n.noticeChannelVersion.Save(txCtx, noticeChannelVersion); err != nil {
			return err
		}

		noticeStat := &notice_entry.NoticeChannelStat{
			NoticeChannelID: noticeChannel.Id,
			VersionID:       noticeChannelVersion.Id,
		}

		if err := n.noticeChannelStat.Save(txCtx, noticeStat); err != nil {
			return err
		}

		n.noticeChannelDriver.RegisterDriver(noticeChannel.Name, noticeChannelDriver)

		return nil
	})
}

func (n *noticeChannelService) UpdateNoticeChannel(ctx context.Context, namespaceId, userID int, channel *notice_model.NoticeChannel) error {
	//1.查看name是否存在
	noticeChannel, err := n.noticeChannelStore.GetByName(ctx, namespaceId, channel.Name)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: noticeChannel.Name,
		Name: noticeChannel.Title,
	})

	t := time.Now()
	noticeChannel.Name = channel.Name
	noticeChannel.Title = channel.Title
	noticeChannel.NamespaceID = namespaceId
	noticeChannel.Type = channel.Type
	noticeChannel.Operator = userID
	noticeChannel.UpdateTime = t

	var noticeChannelDriver notice.IDriverNoticeChannel
	//邮箱
	if channel.Type == 2 {
		email := new(notice_model.NoticeChannelEmail)
		if err = json.Unmarshal([]byte(channel.Config), email); err != nil {
			return err
		}
		noticeChannelDriver = common.NewSmtp(email.SmtpUrl, email.SmtpPort, email.Protocol, email.Account, email.Password, email.Email)
	} else {
		webhook := new(notice_model.NoticeChannelWebhook)
		if err = json.Unmarshal([]byte(channel.Config), webhook); err != nil {
			return err
		}
		noticeChannelDriver = common.NewWebhook(webhook.Url, webhook.Method, webhook.ContentType, webhook.NoticeType, webhook.UserSeparator, webhook.Header, webhook.Template)

	}

	return n.noticeChannelStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = n.noticeChannelStore.Save(txCtx, noticeChannel); err != nil {
			return err
		}

		noticeChannelVersion := &notice_entry.NoticeChannelVersion{
			NoticeChannelID: noticeChannel.Id,
			NamespaceID:     namespaceId,
			NoticeChannelVersionConfig: notice_entry.NoticeChannelVersionConfig{
				Config: channel.Config,
			},
			Operator:   userID,
			CreateTime: t,
		}

		if err = n.noticeChannelVersion.Save(txCtx, noticeChannelVersion); err != nil {
			return err
		}

		noticeStat := &notice_entry.NoticeChannelStat{
			NoticeChannelID: noticeChannel.Id,
			VersionID:       noticeChannelVersion.Id,
		}

		if err = n.noticeChannelStat.Save(txCtx, noticeStat); err != nil {
			return err
		}

		n.noticeChannelDriver.RegisterDriver(noticeChannel.Name, noticeChannelDriver)

		return nil
	})

}

func (n *noticeChannelService) DeleteNoticeChannel(ctx context.Context, namespaceId, userID int, name string) error {
	//1.查看name是否存在
	noticeChannel, err := n.noticeChannelStore.GetByName(ctx, namespaceId, name)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: noticeChannel.Name,
		Name: noticeChannel.Title,
	})

	if !n.isDelete(ctx, noticeChannel.Id) {
		return errors.New("该webhook被告警策略引用，不可删除")
	}

	return n.noticeChannelStore.Transaction(ctx, func(txCtx context.Context) error {

		if _, err = n.noticeChannelStore.Delete(txCtx, noticeChannel.Id); err != nil {
			return err
		}

		delMap := make(map[string]interface{})
		delMap["`kind`"] = "notice_channel"
		delMap["`target`"] = noticeChannel.Id

		if _, err = n.noticeChannelStat.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = n.noticeChannelVersion.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		n.noticeChannelDriver.DelDriver(noticeChannel.Name)
		return nil
	})
}

func (n *noticeChannelService) NoticeChannelList(ctx context.Context, namespaceId int, typ_ int) ([]*notice_model.NoticeChannel, error) {
	noticeChannels, err := n.noticeChannelStore.GetByType(ctx, namespaceId, typ_)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(noticeChannels, func(t *notice_entry.NoticeChannel) int {
		return t.Operator
	})

	infoMaps, _ := n.userService.GetUserInfoMaps(ctx, userIds...)

	channels := make([]*notice_model.NoticeChannel, 0, len(noticeChannels))
	for _, channel := range noticeChannels {
		version, err := n.latestNoticeChannelVersion(ctx, channel.Id)
		if err != nil {
			return nil, err
		}
		operator := ""
		if userInfo, ok := infoMaps[channel.Operator]; ok {
			operator = userInfo.NickName
		}

		channels = append(channels, &notice_model.NoticeChannel{
			Id:         channel.Id,
			Name:       channel.Name,
			Title:      channel.Title,
			Type:       channel.Type,
			Config:     version.Config,
			IsDelete:   n.isDelete(ctx, channel.Id),
			Operator:   operator,
			CreateTime: channel.CreateTime,
			UpdateTime: channel.UpdateTime,
		})
	}

	return channels, nil
}

func (n *noticeChannelService) NoticeChannelByName(ctx context.Context, namespaceId int, name string) (*notice_model.NoticeChannel, error) {
	noticeChannel, err := n.noticeChannelStore.GetByName(ctx, namespaceId, name)
	if err != nil {
		return nil, err
	}
	version, err := n.latestNoticeChannelVersion(ctx, noticeChannel.Id)
	if err != nil {
		return nil, err
	}
	result := &notice_model.NoticeChannel{
		Id:     noticeChannel.Id,
		Name:   noticeChannel.Name,
		Title:  noticeChannel.Title,
		Type:   noticeChannel.Type,
		Config: version.Config,
	}

	return result, nil
}

func (n *noticeChannelService) latestNoticeChannelVersion(ctx context.Context, noticeChannelId int) (*notice_model.NoticeChannelVersion, error) {
	var err error
	stat, err := n.noticeChannelStat.Get(ctx, noticeChannelId)
	if err != nil {
		return nil, err
	}

	var version *notice_entry.NoticeChannelVersion

	version, err = n.noticeChannelVersion.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	return (*notice_model.NoticeChannelVersion)(version), nil
}

func (n *noticeChannelService) isDelete(ctx context.Context, noticeChannelId int) bool {
	quote, _ := n.quoteStore.GetTargetQuote(ctx, noticeChannelId, quote_entry.QuoteTargetKindTypeNoticeChannel)
	return len(quote) == 0
}
