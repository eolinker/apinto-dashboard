package user_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"time"
)

type IUserInfoStore interface {
	store.IBaseStore[user_entry.UserInfo]
	GetByUserName(ctx context.Context, userName string) (*user_entry.UserInfo, error)
	GetAll(ctx context.Context) ([]*user_entry.UserInfo, error)
	UpdateLoginTime(ctx context.Context, userId int, t time.Time) error
	GetUserAmount(ctx context.Context) (int, error)
}

type userInfoStore struct {
	*store.BaseStore[user_entry.UserInfo]
}

func newUserInfoStore(db store.IDB) IUserInfoStore {
	return &userInfoStore{BaseStore: store.CreateStore[user_entry.UserInfo](db)}
}

func (u *userInfoStore) GetByUserName(ctx context.Context, userName string) (*user_entry.UserInfo, error) {
	return u.First(ctx, map[string]interface{}{"username": userName})
}

func (u *userInfoStore) GetAll(ctx context.Context) ([]*user_entry.UserInfo, error) {
	db := u.DB(ctx)

	list := make([]*user_entry.UserInfo, 0)
	err := db.Order("username asc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u *userInfoStore) UpdateLoginTime(ctx context.Context, userId int, t time.Time) error {
	_, err := u.UpdateWhere(ctx, &user_entry.UserInfo{Id: userId}, map[string]interface{}{"login_time": t})
	return err
}

// GetUserAmount 获取未被删除的用户数量
func (u *userInfoStore) GetUserAmount(ctx context.Context) (int, error) {
	db := u.DB(ctx)
	count := int64(0)
	//获取未被删除的用户数量
	if err := db.Model(u.TargetType).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
