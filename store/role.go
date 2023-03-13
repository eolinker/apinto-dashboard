package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IRoleStore interface {
	IBaseStore[entry.Role]
	GetByUUID(ctx context.Context, uuid string) (*entry.Role, error)
	GetByUUIDS(ctx context.Context, uuid []string) ([]*entry.Role, error)
	GetAllRole(ctx context.Context) ([]*entry.Role, error)
	GetByTitle(ctx context.Context, title string) (*entry.Role, error)
}

type roleStore struct {
	*baseStore[entry.Role]
}

func newRoleStore(db IDB) IRoleStore {
	return &roleStore{baseStore: createStore[entry.Role](db)}
}

func (r *roleStore) GetByUUID(ctx context.Context, uuid string) (*entry.Role, error) {
	return r.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (r *roleStore) GetByUUIDS(ctx context.Context, uuid []string) ([]*entry.Role, error) {
	return r.ListQuery(ctx, "`uuid` in (?)", []interface{}{uuid}, "")
}

func (r *roleStore) GetAllRole(ctx context.Context) ([]*entry.Role, error) {
	list := make([]*entry.Role, 0)
	err := r.DB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *roleStore) GetByTitle(ctx context.Context, title string) (*entry.Role, error) {
	return r.First(ctx, map[string]interface{}{"`title`": title})
}
