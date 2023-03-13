package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IWarnStrategyIStore interface {
	IBaseStore[entry.WarnStrategy]
	GetByUuid(ctx context.Context, uuid string) (*entry.WarnStrategy, error)
	GetPage(ctx context.Context, namespaceId, partitionId int, name string, dimension []string, status, pageNum, pageSize int) ([]*entry.WarnStrategy, int64, error)
	GetAll(ctx context.Context, namespaceId, status int) ([]*entry.WarnStrategy, error)
	UpdateIsEnable(ctx context.Context, id int, isEnable bool) error
	GetByTitle(ctx context.Context, namespaceId, partitionId int, title string) (*entry.WarnStrategy, error)
	GetByPartitionId(ctx context.Context, namespaceId, partitionId int) ([]*entry.WarnStrategy, error)
}

type warnStrategyStore struct {
	*BaseStore[entry.WarnStrategy]
}

func newWarnStrategyIStore(db IDB) IWarnStrategyIStore {
	return &warnStrategyStore{BaseStore: CreateStore[entry.WarnStrategy](db)}
}

func (w *warnStrategyStore) GetByUuid(ctx context.Context, uuid string) (*entry.WarnStrategy, error) {
	return w.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (w *warnStrategyStore) GetByTitle(ctx context.Context, namespaceId, partitionId int, title string) (*entry.WarnStrategy, error) {
	return w.FirstQuery(ctx, "`namespace` = ? and `partition_id` = ? and `title` = ?", []interface{}{namespaceId, partitionId, title}, "")
}

func (w *warnStrategyStore) UpdateIsEnable(ctx context.Context, id int, isEnable bool) error {
	_, err := w.UpdateWhere(ctx, &entry.WarnStrategy{Id: id}, map[string]interface{}{"is_enable": isEnable})
	if err != nil {
		return err
	}
	return nil
}

func (w *warnStrategyStore) GetAll(ctx context.Context, namespaceId, status int) ([]*entry.WarnStrategy, error) {
	db := w.DB(ctx).Where("`namespace` = ?", namespaceId)
	if status > -1 {
		db = db.Where("`is_enable` = ?", status)
	}
	list := make([]*entry.WarnStrategy, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (w *warnStrategyStore) GetByPartitionId(ctx context.Context, namespaceId, partitionId int) ([]*entry.WarnStrategy, error) {
	db := w.DB(ctx).Where("`namespace` = ?", namespaceId)
	db = db.Where("`partition_id` = ?", partitionId)
	list := make([]*entry.WarnStrategy, 0)
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (w *warnStrategyStore) GetPage(ctx context.Context, namespaceId, partitionId int, name string, dimension []string, status, pageNum, pageSize int) ([]*entry.WarnStrategy, int64, error) {
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

	list := make([]*entry.WarnStrategy, 0)

	count := int64(0)

	if err := db.Model(list).Count(&count).Limit(pageSize).Offset(entry.PageIndex(pageNum, pageSize)).Order("`update_time` desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}
