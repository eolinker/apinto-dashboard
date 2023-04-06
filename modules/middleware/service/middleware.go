package service

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"gorm.io/gorm"

	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"

	"github.com/go-basic/uuid"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/middleware/entry"

	"github.com/eolinker/apinto-dashboard/modules/middleware"
	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
	"github.com/eolinker/apinto-dashboard/modules/middleware/store"
	"github.com/eolinker/eosc/common/bean"
)

var (
	errPrefixExist             = errors.New("prefix exist")
	errMiddlewareGroupNotFound = errors.New("middleware group not found")
)

type MiddlewareService struct {
	middlewareStore store.IMiddlewareStore
}

func (m *MiddlewareService) CreateGroup(ctx context.Context, namespaceId int, operator int, uuidStr, prefix string, middlewares []string) error {
	_, err := m.middlewareStore.First(ctx, map[string]interface{}{
		"prefix": prefix,
	})
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		// 存在前缀相同的分组
		return errPrefixExist
	}

	if uuidStr == "" {
		uuidStr = uuid.New()
	}
	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuidStr,
		Name: prefix,
	})
	now := time.Now()
	ms, _ := json.Marshal(middlewares)

	return m.middlewareStore.Insert(ctx, &entry.Middleware{
		NamespaceId: namespaceId,
		UUID:        uuidStr,
		Prefix:      prefix,
		Middlewares: string(ms),
		Operator:    operator,
		CreateTime:  now,
		UpdateTime:  now,
	})
}

func (m *MiddlewareService) UpdateGroup(ctx context.Context, namespaceId int, operator int, uuidStr, prefix string, middlewares []string) error {

	info, err := m.middlewareStore.First(ctx, map[string]interface{}{
		"prefix": prefix,
	})
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		info, err = m.middlewareStore.GetByUUID(ctx, uuidStr)
		if err != nil {
			return err
		}
	} else if info.UUID != uuidStr {
		return errPrefixExist
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuidStr,
		Name: prefix,
	})
	ms, _ := json.Marshal(middlewares)

	info.Prefix = prefix
	info.Operator = operator
	info.Middlewares = string(ms)
	info.UpdateTime = time.Now()
	return m.middlewareStore.Save(ctx, info)
}

func (m *MiddlewareService) DeleteGroup(ctx context.Context, namespaceId int, operator int, uuidStr string) error {
	info, err := m.middlewareStore.GetByUUID(ctx, uuidStr)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		return nil
	}
	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuidStr,
		Name: info.Prefix,
	})
	_, err = m.middlewareStore.Delete(ctx, info.Id)
	return err
}

func (m *MiddlewareService) GroupList(ctx context.Context, namespaceId int) ([]*model.MiddlewareGroup, error) {
	list, err := m.middlewareStore.List(ctx, map[string]interface{}{
		"namespace_id": namespaceId,
	})
	if err != nil {
		return nil, err
	}
	groups := make([]*model.MiddlewareGroup, 0, len(list))
	for _, g := range list {
		groups = append(groups, &model.MiddlewareGroup{
			Uuid:   g.UUID,
			Prefix: g.Prefix,
		})
	}
	sort.Sort(model.MiddlewareGroups(groups))
	return groups, nil
}

func (m *MiddlewareService) GroupInfo(ctx context.Context, namespaceId int, uuidStr string) (*model.MiddlewareGroupInfo, error) {
	info, err := m.middlewareStore.GetByUUID(ctx, uuidStr)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, errMiddlewareGroupNotFound
	}
	var middlewares = []string{}
	err = json.Unmarshal([]byte(info.Middlewares), &middlewares)
	if err != nil {
		return nil, err
	}
	return &model.MiddlewareGroupInfo{
		Middlewares: middlewares,
		MiddlewareGroup: model.MiddlewareGroup{
			Uuid:   info.UUID,
			Prefix: info.UUID,
		},
		All: demoMiddlewares,
	}, nil
}

func newMiddlewareService() middleware.IMiddlewareService {
	c := &MiddlewareService{}
	bean.Autowired(&c.middlewareStore)
	return c
}
