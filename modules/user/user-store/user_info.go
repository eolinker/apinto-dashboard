package user_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"time"
)

type IUserInfoStore interface {
	store.IBaseStore[user_entry.UserInfo]
	GetByUUID(ctx context.Context, uuid string) (*user_entry.UserInfo, error)
	GetByUserName(ctx context.Context, userName string) (*user_entry.UserInfo, error)
	GetByRoleByKeyword(ctx context.Context, roleId int, keyword string) ([]*user_entry.UserInfo, error)
	GetAll(ctx context.Context) ([]*user_entry.UserInfo, error)
	UpdateStatus(ctx context.Context, userId, status int) error
	UpdateLoginTime(ctx context.Context, userId int, t time.Time) error
	SoftDelete(ctx context.Context, operator int, userId ...int) error
	GetUserAmount(ctx context.Context) (int, error)
}

type userInfoStore struct {
	*store.BaseStore[user_entry.UserInfo]
}

func newUserInfoStore(db store.IDB) IUserInfoStore {
	return &userInfoStore{BaseStore: store.CreateStore[user_entry.UserInfo](db)}
}

func (u *userInfoStore) GetByUUID(ctx context.Context, uuid string) (*user_entry.UserInfo, error) {
	return u.First(ctx, map[string]interface{}{"uuid": uuid})
}

func (u *userInfoStore) GetByUserName(ctx context.Context, userName string) (*user_entry.UserInfo, error) {
	return u.First(ctx, map[string]interface{}{"user_name": userName, "is_delete": 0})
}

func (u *userInfoStore) GetByRoleByKeyword(ctx context.Context, roleId int, keyword string) ([]*user_entry.UserInfo, error) {
	db := u.DB(ctx).Where("is_delete = 0")

	if roleId != 0 {
		db = db.Where("`id` in (select user_id from user_role where role_id = ?)", roleId)
	}

	if keyword != "" {
		keyword = "%" + keyword + "%"
		db = db.Where("`user_name` like ? or `nick_name` like ?", keyword, keyword)
	}

	list := make([]*user_entry.UserInfo, 0)
	err := db.Order("update_time desc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u *userInfoStore) GetAll(ctx context.Context) ([]*user_entry.UserInfo, error) {
	db := u.DB(ctx).Where("is_delete = 0")

	list := make([]*user_entry.UserInfo, 0)
	err := db.Order("update_time desc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u *userInfoStore) UpdateStatus(ctx context.Context, userId, status int) error {
	_, err := u.UpdateWhere(ctx, &user_entry.UserInfo{Id: userId}, map[string]interface{}{"status": status})
	return err
}

func (u *userInfoStore) UpdateLoginTime(ctx context.Context, userId int, t time.Time) error {
	_, err := u.UpdateWhere(ctx, &user_entry.UserInfo{Id: userId}, map[string]interface{}{"last_login_time": t})
	return err
}

func (u *userInfoStore) SoftDelete(ctx context.Context, operator int, userIds ...int) error {

	for _, id := range userIds {
		_, err := u.UpdateWhere(ctx, &user_entry.UserInfo{Id: id}, map[string]interface{}{"is_delete": true, "operator": operator})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserAmount 获取未被删除的用户数量
func (u *userInfoStore) GetUserAmount(ctx context.Context) (int, error) {
	db := u.DB(ctx)
	count := int64(0)
	//获取未被删除的用户数量
	if err := db.Model(u.TargetType).Where("`is_delete` = 0 ").Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
