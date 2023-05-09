package user

import (
	"context"
	"time"

	"github.com/eolinker/apinto-dashboard/cache"
	user_dto "github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	"github.com/eolinker/apinto-dashboard/modules/user/user-model"
)

type IUserInfoService interface {
	GetAllUsers(ctx context.Context) ([]*user_model.UserInfo, error)
	GetUserInfo(ctx context.Context, userId int) (*user_model.UserInfo, error)
	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*user_model.UserInfo, error)
	CheckPassword(ctx context.Context, name string, password string) (int, bool)
	GetUserInfoByName(ctx context.Context, userName string) (*user_model.UserInfo, error)
	GetUserInfoByNames(ctx context.Context, userNames ...string) (map[string]*user_model.UserInfo, error)
	UpdateMyProfile(ctx context.Context, userId int, req *user_dto.UpdateMyProfileReq) error
	UpdateMyPassword(ctx context.Context, userId int, req *user_dto.UpdateMyPasswordReq) error
	UpdateLastLoginTime(ctx context.Context, userId int, loginTime *time.Time) error
}

type ISessionCache interface {
	cache.IRedisCache[user_model.Session, string]
}
