package service

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/user"
	user_entry "github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	user_store "github.com/eolinker/apinto-dashboard/modules/user/user-store"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

const AdminName = "admin"
const defaultPwd = "12345678"

type userInfoService struct {
	userInfoStore user_store.IUserInfoStore
	cache         cache.IRedisCache[user_model.UserInfo]
}

func newUserInfoService() user.IUserInfoService {
	u := &userInfoService{}
	bean.Autowired(&u.userInfoStore)
	return u
}

func (u *userInfoService) GetUserInfoMaps(ctx context.Context, userIds ...int) (map[int]*user_model.UserInfo, error) {

	maps := make(map[int]*user_model.UserInfo)
	need := make([]int, 0, len(userIds))

	for _, userId := range userIds {
		userInfo, ok := u.getCache(ctx, userId)
		if ok {
			maps[userInfo.Id] = userInfo
		} else {
			need = append(need, userId)
		}
	}
	if len(need) > 0 {
		if len(need) < 5 {
			for _, userId := range need {
				userInfo, err := u.GetUserInfo(ctx, userId)
				if err == nil {
					maps[userInfo.Id] = userInfo
				} else {

				}
			}
		}
	}

	return maps, nil
}
func (u *userInfoService) getCache(ctx context.Context, userID int) (*user_model.UserInfo, bool) {
	key := fmt.Sprintf("apinto:userinfo-id:%d", userID)
	userModel, err := u.cache.Get(ctx, key)
	if err != nil {
		return nil, false
	}
	return userModel, true
}

// GetUserInfo 获取不到用户信息记录错误即可，不必返回error
func (u *userInfoService) GetUserInfo(ctx context.Context, userID int) (*user_model.UserInfo, error) {
	key := fmt.Sprintf("apinto:userinfo-id:%d", userID)
	userModel, err := u.cache.Get(ctx, key)
	if err == nil {
		return userModel, nil
	}
	userInfo, err := u.userInfoStore.Get(ctx, userID)
	if err != nil {
		userModel = &user_model.UserInfo{
			Id:            userID,
			Sex:           0,
			UserName:      "unknown",
			NoticeUserId:  "",
			NickName:      "unknown",
			Email:         "unknown",
			Phone:         "unknown",
			Avatar:        "",
			LastLoginTime: nil,
		}
	} else {
		userModel = entryToModule(userInfo)
	}
	u.cache.Set(ctx, key, userModel, time.Hour)
	return userModel, nil
}

func entryToModule(info *user_entry.UserInfo) *user_model.UserInfo {

	return &user_model.UserInfo{
		Id:            info.Id,
		Sex:           info.Sex,
		UserName:      info.UserName,
		NoticeUserId:  info.NoticeUserId,
		NickName:      info.NickName,
		Email:         info.Email,
		Phone:         info.Phone,
		Avatar:        info.Avatar,
		LastLoginTime: info.LastLoginTime,
	}
}
