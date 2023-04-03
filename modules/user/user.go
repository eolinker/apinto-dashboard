package user

import (
	"context"
	"github.com/eolinker/apinto-dashboard/access"
	user_model2 "github.com/eolinker/apinto-dashboard/modules/user/user-model"
)

type IUserInfoService interface {
	GetUserInfo(ctx context.Context, userId int) (*user_model2.UserInfo, error)

	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*user_model2.UserInfo, error)

	GetAccessInfo(ctx context.Context, userId int) (map[access.Access]struct{}, error)
}
