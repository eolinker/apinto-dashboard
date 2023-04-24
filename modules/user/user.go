package user

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/user/user-model"
)

type IUserInfoService interface {
	GetUserInfo(ctx context.Context, userId int) (*user_model.UserInfo, error)
	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*user_model.UserInfo, error)
	GetUserInfoByName(ctx context.Context, userName string) (*user_model.UserInfo, error)
	GetUserInfoByNames(ctx context.Context, userNames ...string) (map[string]*user_model.UserInfo, error)
}
