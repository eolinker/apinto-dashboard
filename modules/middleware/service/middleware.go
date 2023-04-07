package service

import (
	"context"
	"encoding/json"
	system_entry "github.com/eolinker/apinto-dashboard/modules/base/system-entry"
	system_store "github.com/eolinker/apinto-dashboard/modules/base/system-store"

	"github.com/eolinker/apinto-dashboard/modules/middleware"

	"gorm.io/gorm"

	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
)

var (
	systemKey = "middleware"
)

func newMiddlewareService() middleware.IMiddlewareService {
	c := &MiddlewareService{}
	bean.Autowired(&c.systemStore)
	return c
}

type MiddlewareService struct {
	systemStore system_store.ISystemInfoStore
}

func (m *MiddlewareService) Save(ctx context.Context, config []*model.MiddlewareGroup) error {

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	value, err := m.systemStore.GetSystemInfoByKey(ctx, systemKey)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		value = &system_entry.SystemInfo{
			Key: systemKey,
		}
	}
	value.Value = data
	return m.systemStore.Save(ctx, value)
}

func (m *MiddlewareService) Groups(ctx context.Context) ([]*model.MiddlewareGroup, error) {
	value, err := m.systemStore.GetSystemInfoByKey(ctx, systemKey)
	groups := make([]*model.MiddlewareGroup, 0)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(value.Value), &groups)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil

}
