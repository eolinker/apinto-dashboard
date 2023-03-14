package notice

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-model"
)

type INoticeChannelService interface {
	CreateNoticeChannel(ctx context.Context, namespaceId, userID int, channel *notice_model.NoticeChannel) error
	UpdateNoticeChannel(ctx context.Context, namespaceId, userID int, channel *notice_model.NoticeChannel) error
	DeleteNoticeChannel(ctx context.Context, namespaceId, userID int, name string) error
	NoticeChannelList(ctx context.Context, namespaceId, typ_ int) ([]*notice_model.NoticeChannel, error)
	NoticeChannelByName(ctx context.Context, namespaceId int, name string) (*notice_model.NoticeChannel, error)
	InitChannelDriver() error
}
