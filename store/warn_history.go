package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"time"
)

type IWarnHistoryIStore interface {
	IBaseStore[entry.WarnHistory]
	GetPage(ctx context.Context, partitionId int, name string, pageNum, pageSize int, startTime, endTime time.Time) ([]*entry.WarnHistory, int64, error)
}

type WarnHistoryStore struct {
	*BaseStore[entry.WarnHistory]
}

func newWarnHistoryIStore(db IDB) IWarnHistoryIStore {
	return &WarnHistoryStore{BaseStore: CreateStore[entry.WarnHistory](db)}
}

func (w *WarnHistoryStore) GetPage(ctx context.Context, partitionId int, name string, pageNum, pageSize int, startTime, endTime time.Time) ([]*entry.WarnHistory, int64, error) {
	db := w.DB(ctx).Where("`partition_id` = ?", partitionId)
	if name != "" {
		db = db.Where("`strategy_title` like ?", "%"+name+"%")
	}
	db = db.Where("`create_time` >= ? and `create_time` <= ?", startTime, endTime)

	list := make([]*entry.WarnHistory, 0)

	count := int64(0)

	if err := db.Model(list).Count(&count).Limit(pageSize).Offset(entry.PageIndex(pageNum, pageSize)).Order("`create_time` desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
