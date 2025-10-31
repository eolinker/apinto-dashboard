package user

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/user/user-model"
)

type UserInfoIterate interface {
	UserId() int
	Set(name string)
}

type IUserInfoService interface {
	GetAllUsers(ctx context.Context) ([]*user_model.UserInfo, error)
	GetUserInfo(ctx context.Context, userId int) (*user_model.UserInfo, error)
	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*user_model.UserInfo, error)
	//CheckPassword(ctx context.Context, name string, password string) (int, bool)
	GetUserInfoByName(ctx context.Context, userName string) (*user_model.UserInfo, error)
	GetUserInfoByNames(ctx context.Context, userNames ...string) (map[string]*user_model.UserInfo, error)
	SaveUserInfo(ctx context.Context, info *user_model.UserInfo) error
	SetUserName(ctx context.Context, iterate ...UserInfoIterate) error
	//UpdateMyProfile(ctx context.Context, userId int, req *user_dto.UpdateMyProfileReq) error
	//UpdateMyPassword(ctx context.Context, userId int, req *user_dto.UpdateMyPasswordReq) error
	//UpdateLastLoginTime(ctx context.Context, userId int, loginTime *time.Time) error

}

func SetUserName[T UserInfoIterate](s IUserInfoService, ctx context.Context, ts ...T) error {
	if len(ts) == 0 {
		return nil
	}
	tl := make([]UserInfoIterate, 0, len(ts))
	for _, i := range ts {
		tl = append(tl, i)
	}
	return s.SetUserName(ctx, tl...)
}

//type ISessionCache interface {
//	cache.IRedisCache[user_model.Session, string]
//}
