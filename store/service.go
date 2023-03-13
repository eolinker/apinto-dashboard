package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IServiceStore interface {
	IBaseStore[entry.Service]
	GetListPage(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*entry.Service, int, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*entry.Service, error)
	GetListAll(ctx context.Context, namespaceId int) ([]*entry.Service, error)
	GetByNames(ctx context.Context, namespaceId int, names []string) ([]*entry.Service, error)
}

type serviceStore struct {
	*baseStore[entry.Service]
}

func newServiceStore(db IDB) IServiceStore {
	return &serviceStore{baseStore: createStore[entry.Service](db)}
}

func (s *serviceStore) GetListPage(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*entry.Service, int, error) {
	services := make([]*entry.Service, 0)
	db := s.DB(ctx).Where("`namespace` = ?", namespaceID)
	count := int64(0)
	var err error
	if searchName != "" {
		db = db.Where("`name` like ?", "%"+searchName+"%")
	}
	if pageNum > 0 && pageSize > 0 {
		err = db.Model(services).Count(&count).Order("`update_time` DESC").Limit(pageSize).Offset(entry.PageIndex(pageNum, pageSize)).Find(&services).Error
	} else {
		err = db.Model(services).Count(&count).Order("`update_time` DESC").Find(&services).Error
	}
	if err != nil {
		return nil, 0, err
	}

	return services, int(count), nil
}

func (s *serviceStore) GetByName(ctx context.Context, namespaceId int, name string) (*entry.Service, error) {
	return s.FirstQuery(ctx, "`namespace` = ? and `name` = ?", []interface{}{namespaceId, name}, "")
}

func (s *serviceStore) GetByUUIDs(ctx context.Context, namespaceId int, uuids []string) ([]*entry.Service, error) {
	services := make([]*entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ? and `name` in (?)", namespaceId, uuids).Find(&services).Error
	return services, err
}

func (s *serviceStore) GetByNames(ctx context.Context, namespaceId int, names []string) ([]*entry.Service, error) {
	services := make([]*entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ? and `name` in (?)", namespaceId, names).Find(&services).Error
	return services, err
}

func (s *serviceStore) GetListAll(ctx context.Context, namespaceId int) ([]*entry.Service, error) {
	services := make([]*entry.Service, 0)
	err := s.DB(ctx).Where("`namespace` = ?", namespaceId).Find(&services).Error
	return services, err
}
