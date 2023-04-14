package upstream_store

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/upstream/upstream-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceStore interface {
	store.IBaseStore[upstream_entry.Service]
	GetListPage(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*upstream_entry.Service, int, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*upstream_entry.Service, error)
	GetListAll(ctx context.Context, namespaceId int) ([]*upstream_entry.Service, error)
	GetByNames(ctx context.Context, namespaceId int, names []string) ([]*upstream_entry.Service, error)
	ServiceCount(ctx context.Context, params map[string]interface{}) (int64, error)
}

type serviceStore struct {
	*store.BaseStore[upstream_entry.Service]
}

func newServiceStore(db store.IDB) IServiceStore {
	return &serviceStore{BaseStore: store.CreateStore[upstream_entry.Service](db)}
}

func (s *serviceStore) ServiceCount(ctx context.Context, params map[string]interface{}) (int64, error) {
	var count int64
	err := s.DB(ctx).Where(params).Model(upstream_entry.Service{}).Count(&count).Error
	return count, err
}

func (s *serviceStore) GetListPage(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*upstream_entry.Service, int, error) {
	services := make([]*upstream_entry.Service, 0)
	db := s.DB(ctx).Where("`namespace` = ?", namespaceID)
	count := int64(0)
	var err error
	if searchName != "" {
		db = db.Where("`name` like ?", "%"+searchName+"%")
	}
	if pageNum > 0 && pageSize > 0 {
		err = db.Model(services).Count(&count).Order("`update_time` DESC").Limit(pageSize).Offset(store.PageIndex(pageNum, pageSize)).Find(&services).Error
	} else {
		err = db.Model(services).Count(&count).Order("`update_time` DESC").Find(&services).Error
	}
	if err != nil {
		return nil, 0, err
	}

	return services, int(count), nil
}

func (s *serviceStore) GetByName(ctx context.Context, namespaceId int, name string) (*upstream_entry.Service, error) {
	return s.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (s *serviceStore) GetByUUIDs(ctx context.Context, namespaceId int, uuids []string) ([]*upstream_entry.Service, error) {
	services := make([]*upstream_entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ? and `name` in (?)", namespaceId, uuids).Find(&services).Error
	return services, err
}

func (s *serviceStore) GetByNames(ctx context.Context, namespaceId int, names []string) ([]*upstream_entry.Service, error) {
	services := make([]*upstream_entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ? and `name` in (?)", namespaceId, names).Find(&services).Error
	return services, err
}

func (s *serviceStore) GetListAll(ctx context.Context, namespaceId int) ([]*upstream_entry.Service, error) {
	services := make([]*upstream_entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ?", namespaceId).Find(&services).Error
	return services, err
}
