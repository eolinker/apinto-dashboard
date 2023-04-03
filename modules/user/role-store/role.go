package role_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/user/role-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IRoleStore interface {
	store.IBaseStore[role_entry.Role]
	GetByUUID(ctx context.Context, uuid string) (*role_entry.Role, error)
	GetByUUIDS(ctx context.Context, uuid []string) ([]*role_entry.Role, error)
	GetAllRole(ctx context.Context) ([]*role_entry.Role, error)
	GetByTitle(ctx context.Context, title string) (*role_entry.Role, error)
}

type roleStore struct {
	*store.BaseStore[role_entry.Role]
}

func newRoleStore(db store.IDB) IRoleStore {
	return &roleStore{BaseStore: store.CreateStore[role_entry.Role](db)}
}

func (r *roleStore) GetByUUID(ctx context.Context, uuid string) (*role_entry.Role, error) {
	return r.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (r *roleStore) GetByUUIDS(ctx context.Context, uuid []string) ([]*role_entry.Role, error) {
	return r.ListQuery(ctx, "`uuid` in (?)", []interface{}{uuid}, "")
}

func (r *roleStore) GetAllRole(ctx context.Context) ([]*role_entry.Role, error) {
	list := make([]*role_entry.Role, 0)
	err := r.DB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *roleStore) GetByTitle(ctx context.Context, title string) (*role_entry.Role, error) {
	return r.First(ctx, map[string]interface{}{"`title`": title})
}
