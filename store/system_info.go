package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"gorm.io/gorm"
)

var (
	_ ISystemInfoStore = (*systemInfoStore)(nil)
)

type ISystemInfoStore interface {
	IBaseStore[entry.SystemInfo]
	GetSystemInfoByKey(ctx context.Context, key string) (*entry.SystemInfo, error)
	InitDashboardID(ctx context.Context, id string) error
}

type systemInfoStore struct {
	*BaseStore[entry.SystemInfo]
}

func newSystemInfoStore(db IDB) ISystemInfoStore {
	return &systemInfoStore{BaseStore: CreateStore[entry.SystemInfo](db)}
}

func (s *systemInfoStore) GetSystemInfoByKey(ctx context.Context, key string) (*entry.SystemInfo, error) {
	config := new(entry.SystemInfo)
	err := s.DB(ctx).Where("`key` = ?", key).Take(config).Error
	return config, err
}

func (s *systemInfoStore) InitDashboardID(ctx context.Context, id string) error {
	_, err := s.FirstQuery(ctx, "`key` = ?", []interface{}{enum.DashboardIdDBKey}, "")
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		err = s.Save(ctx, &entry.SystemInfo{
			Key:   enum.DashboardIdDBKey,
			Value: []byte(id),
		})
	}
	return err
}
