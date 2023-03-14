package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/monitor-entry"
	"github.com/eolinker/apinto-dashboard/entry/page-entry"
	"time"
)

type IWarnHistoryIStore interface {
	IBaseStore[monitor_entry.WarnHistory]
	GetPage(ctx context.Context, partitionId int, name string, pageNum, pageSize int, startTime, endTime time.Time) ([]*monitor_entry.WarnHistory, int64, error)
}

type WarnHistoryStore struct {
	*BaseStore[monitor_entry.WarnHistory]
}

func newWarnHistoryIStore(db IDB) IWarnHistoryIStore {
	return &WarnHistoryStore{BaseStore: CreateStore[monitor_entry.WarnHistory](db)}
}

func (w *WarnHistoryStore) GetPage(ctx context.Context, partitionId int, name string, pageNum, pageSize int, startTime, endTime time.Time) ([]*monitor_entry.WarnHistory, int64, error) {
	db := w.DB(ctx).Where("`partition_id` = ?", partitionId)
	if name != "" {
		db = db.Where("`strategy_title` like ?", "%"+name+"%")
	}
	db = db.Where("`create_time` >= ? and `create_time` <= ?", startTime, endTime)

	list := make([]*monitor_entry.WarnHistory, 0)

	count := int64(0)

	if err := db.Model(list).Count(&count).Limit(pageSize).Offset(page_entry.PageIndex(pageNum, pageSize)).Order("`create_time` desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
