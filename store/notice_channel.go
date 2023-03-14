package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/notice-entry"
)

type INoticeChannelStore interface {
	IBaseStore[notice_entry.NoticeChannel]
	GetByType(ctx context.Context, namespaceId, typ_ int) ([]*notice_entry.NoticeChannel, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*notice_entry.NoticeChannel, error)
	GetAll(ctx context.Context) ([]*notice_entry.NoticeChannel, error)
}

type noticeChannelStore struct {
	*BaseStore[notice_entry.NoticeChannel]
}

func newNoticeChannelStore(db IDB) INoticeChannelStore {
	return &noticeChannelStore{BaseStore: CreateStore[notice_entry.NoticeChannel](db)}
}

func (n *noticeChannelStore) GetByType(ctx context.Context, namespaceId, typ_ int) ([]*notice_entry.NoticeChannel, error) {
	db := n.DB(ctx).Where("`namespace` = ?", namespaceId)
	if typ_ > 0 {
		db = db.Where("`type` = ?", typ_)
	}

	list := make([]*notice_entry.NoticeChannel, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (n *noticeChannelStore) GetAll(ctx context.Context) ([]*notice_entry.NoticeChannel, error) {
	db := n.DB(ctx)

	list := make([]*notice_entry.NoticeChannel, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (n *noticeChannelStore) GetByName(ctx context.Context, namespaceId int, name string) (*notice_entry.NoticeChannel, error) {
	return n.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}
