package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IRoleAccessStore interface {
	IBaseStore[entry.RoleAccess]
	GetByRoleID(ctx context.Context, roleID int) (*entry.RoleAccess, error)
}

type roleAccessStore struct {
	*BaseStore[entry.RoleAccess]
}

func newRoleAccessStore(db IDB) IRoleAccessStore {
	return &roleAccessStore{BaseStore: CreateStore[entry.RoleAccess](db)}
}

func (r *roleAccessStore) GetByRoleID(ctx context.Context, roleID int) (*entry.RoleAccess, error) {
	roleAccess := new(entry.RoleAccess)
	err := r.DB(ctx).Where("`role_id` = ?", roleID).Take(roleAccess).Error
	return roleAccess, err
}
