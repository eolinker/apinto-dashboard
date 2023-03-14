package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/user-entry"
)

type IUserRoleStore interface {
	IBaseStore[user_entry.UserRole]
	GetByUUID(ctx context.Context, uuid string) (*user_entry.UserRole, error)
	DelByUserId(ctx context.Context, userId int) error
	GetRoleQuotedCount(ctx context.Context, roleID int) (int, error)
	GetListByRoleID(ctx context.Context, roleID int) ([]*user_entry.UserRole, error)
	GetByUserIDRoleID(ctx context.Context, userID, roleID int) (*user_entry.UserRole, error)
}

type userRoleStore struct {
	*BaseStore[user_entry.UserRole]
}

func newUserRoleStore(db IDB) IUserRoleStore {
	return &userRoleStore{BaseStore: CreateStore[user_entry.UserRole](db)}
}

func (u *userRoleStore) DelByUserId(ctx context.Context, userId int) error {
	_, err := u.DeleteWhere(ctx, map[string]interface{}{"user_id": userId})
	return err
}

func (u *userRoleStore) GetByUUID(ctx context.Context, uuid string) (*user_entry.UserRole, error) {
	return u.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (u *userRoleStore) GetRoleQuotedCount(ctx context.Context, roleID int) (int, error) {
	db := u.DB(ctx)
	count := int64(0)
	if err := db.Model(u.targetType).Where("`role_id` = ?", roleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (u *userRoleStore) GetListByRoleID(ctx context.Context, roleID int) ([]*user_entry.UserRole, error) {
	list := make([]*user_entry.UserRole, 0)
	err := u.DB(ctx).Where("`role_id` = ?", roleID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *userRoleStore) GetByUserIDRoleID(ctx context.Context, userID, roleID int) (*user_entry.UserRole, error) {
	return u.FirstQuery(ctx, "`user_id` = ? and `role_id` = ?", []interface{}{userID, roleID}, "")
}
