package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/modules/user"
	user_entry "github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"time"
)

const AdminName = "admin"
const defaultPwd = "12345678"

var (
	accessSet = make(map[access.Access]struct{})
)

func init() {
	all := access.All()
	for _, a := range all {
		accessSet[a] = struct{}{}
	}
}

type userInfoService struct {
}

func newUserInfoService() user.IUserInfoService {
	u := &userInfoService{}

	return u
}

func (u *userInfoService) GetUserInfoMaps(ctx context.Context, userIds ...int) (map[int]*user_model.UserInfo, error) {

	maps := make(map[int]*user_model.UserInfo)
	for _, userId := range userIds {
		userInfo, err := u.GetUserInfo(ctx, userId)
		if err != nil {
			continue
		}
		maps[userInfo.Id] = userInfo
	}

	return maps, nil
}

// GetUserInfo 获取不到用户信息记录错误即可，不必返回error
func (u *userInfoService) GetUserInfo(ctx context.Context, userID int) (*user_model.UserInfo, error) {
	return &user_model.UserInfo{
		UserInfo: &user_entry.UserInfo{
			Id:            userID,
			Sex:           0,
			UserName:      AdminName,
			NoticeUserId:  "",
			NickName:      AdminName,
			Email:         "",
			Phone:         "",
			Avatar:        "",
			Remark:        "",
			RoleIds:       "",
			Status:        0,
			IsDelete:      false,
			Operator:      0,
			FlushTime:     time.Time{},
			CreateTime:    time.Time{},
			UpdateTime:    time.Time{},
			LastLoginTime: nil,
		},
		OperateEnable: false,
		Operator:      "",
	}, nil

}

func (u *userInfoService) GetAccessInfo(ctx context.Context, userId int) (map[access.Access]struct{}, error) {
	// todo 临时实现，后续插件化
	return accessSet, nil
}
