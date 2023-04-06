package service

import (
	"context"
	"encoding/json"

	"github.com/eolinker/apinto-dashboard/modules/middleware"

	"gorm.io/gorm"

	"github.com/eolinker/apinto-dashboard/modules/system/entry"
	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/apinto-dashboard/modules/system/store"

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
	systemStore store.ISystemStore
}

func (m *MiddlewareService) Save(ctx context.Context, config string) error {
	value, err := m.systemStore.First(ctx, map[string]interface{}{
		"key": systemKey,
	})
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		value = &entry.System{
			Key: systemKey,
		}
	}
	value.Value = config
	return m.systemStore.Save(ctx, value)
}

func (m *MiddlewareService) Groups(ctx context.Context) (*model.Middleware, error) {
	value, err := m.systemStore.First(ctx, map[string]interface{}{
		"key": systemKey,
	})
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

	return &model.Middleware{
		Group:       groups,
		Middlewares: demoMiddlewares,
	}, nil
}
