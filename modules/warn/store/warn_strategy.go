package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IWarnStrategyIStore interface {
	store.IBaseStore[warn_entry.WarnStrategy]
	GetByUuid(ctx context.Context, uuid string) (*warn_entry.WarnStrategy, error)
	GetPage(ctx context.Context, namespaceId, partitionId int, name string, dimension []string, status, pageNum, pageSize int) ([]*warn_entry.WarnStrategy, int64, error)
	GetAll(ctx context.Context, namespaceId, status int) ([]*warn_entry.WarnStrategy, error)
	UpdateIsEnable(ctx context.Context, id int, isEnable bool) error
	GetByTitle(ctx context.Context, namespaceId, partitionId int, title string) (*warn_entry.WarnStrategy, error)
	GetByPartitionId(ctx context.Context, namespaceId, partitionId int) ([]*warn_entry.WarnStrategy, error)
}

type warnStrategyStore struct {
	*store.BaseStore[warn_entry.WarnStrategy]
}

func newWarnStrategyIStore(db store.IDB) IWarnStrategyIStore {
	return &warnStrategyStore{BaseStore: store.CreateStore[warn_entry.WarnStrategy](db)}
}

func (w *warnStrategyStore) GetByUuid(ctx context.Context, uuid string) (*warn_entry.WarnStrategy, error) {
	return w.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (w *warnStrategyStore) GetByTitle(ctx context.Context, namespaceId, partitionId int, title string) (*warn_entry.WarnStrategy, error) {
	return w.FirstQuery(ctx, "`namespace` = ? and `partition_id` = ? and `title` = ?", []interface{}{namespaceId, partitionId, title}, "")
}

func (w *warnStrategyStore) UpdateIsEnable(ctx context.Context, id int, isEnable bool) error {
	_, err := w.UpdateWhere(ctx, &warn_entry.WarnStrategy{Id: id}, map[string]interface{}{"is_enable": isEnable})
	if err != nil {
		return err
	}
	return nil
}

func (w *warnStrategyStore) GetAll(ctx context.Context, namespaceId, status int) ([]*warn_entry.WarnStrategy, error) {
	db := w.DB(ctx).Where("`namespace` = ?", namespaceId)
	if status > -1 {
		db = db.Where("`is_enable` = ?", status)
	}
	list := make([]*warn_entry.WarnStrategy, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (w *warnStrategyStore) GetByPartitionId(ctx context.Context, namespaceId, partitionId int) ([]*warn_entry.WarnStrategy, error) {
	db := w.DB(ctx).Where("`namespace` = ?", namespaceId)
	db = db.Where("`partition_id` = ?", partitionId)
	list := make([]*warn_entry.WarnStrategy, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (w *warnStrategyStore) GetPage(ctx context.Context, namespaceId, partitionId int, name string, dimension []string, status, pageNum, pageSize int) ([]*warn_entry.WarnStrategy, int64, error) {
	db := w.DB(ctx).Where("`namespace` = ? and `partition_id` = ?", namespaceId, partitionId)
	if name != "" {
		db = db.Where("`title` like ?", "%"+name+"%")
	}
	if len(dimension) > 0 {
		db = db.Where("`dimension` in (?)", dimension)
	}
	if status > -1 {
		db = db.Where("`is_enable` = ?", status)
	}

	list := make([]*warn_entry.WarnStrategy, 0)

	count := int64(0)

	if err := db.Model(list).Count(&count).Limit(pageSize).Offset(store.PageIndex(pageNum, pageSize)).Order("`update_time` desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
