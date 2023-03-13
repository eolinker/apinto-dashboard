package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type INoticeChannelStore interface {
	IBaseStore[entry.NoticeChannel]
	GetByType(ctx context.Context, namespaceId, typ_ int) ([]*entry.NoticeChannel, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*entry.NoticeChannel, error)
	GetAll(ctx context.Context) ([]*entry.NoticeChannel, error)
}

type noticeChannelStore struct {
	*baseStore[entry.NoticeChannel]
}

func newNoticeChannelStore(db IDB) INoticeChannelStore {
	return &noticeChannelStore{baseStore: createStore[entry.NoticeChannel](db)}
}

func (n *noticeChannelStore) GetByType(ctx context.Context, namespaceId, typ_ int) ([]*entry.NoticeChannel, error) {
	db := n.DB(ctx).Where("`namespace` = ?", namespaceId)
	if typ_ > 0 {
		db = db.Where("`type` = ?", typ_)
	}

	list := make([]*entry.NoticeChannel, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (n *noticeChannelStore) GetAll(ctx context.Context) ([]*entry.NoticeChannel, error) {
	db := n.DB(ctx)

	list := make([]*entry.NoticeChannel, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (n *noticeChannelStore) GetByName(ctx context.Context, namespaceId int, name string) (*entry.NoticeChannel, error) {
	return n.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}
