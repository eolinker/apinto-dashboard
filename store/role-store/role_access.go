package role_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/role-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IRoleAccessStore interface {
	store.IBaseStore[role_entry.RoleAccess]
	GetByRoleID(ctx context.Context, roleID int) (*role_entry.RoleAccess, error)
}

type roleAccessStore struct {
	*store.BaseStore[role_entry.RoleAccess]
}

func newRoleAccessStore(db store.IDB) IRoleAccessStore {
	return &roleAccessStore{BaseStore: store.CreateStore[role_entry.RoleAccess](db)}
}

func (r *roleAccessStore) GetByRoleID(ctx context.Context, roleID int) (*role_entry.RoleAccess, error) {
	roleAccess := new(role_entry.RoleAccess)
	err := r.DB(ctx).Where("`role_id` = ?", roleID).Take(roleAccess).Error
	return roleAccess, err
}
