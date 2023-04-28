package service

import (
	"context"
	"errors"
	"time"

	"github.com/eolinker/apinto-dashboard/common"

	"github.com/eolinker/apinto-dashboard/modules/user"
	user_dto "github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	user_entry "github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	user_store "github.com/eolinker/apinto-dashboard/modules/user/user-store"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
)

const AdminName = "admin"
const defaultPwd = "12345678"

type userInfoService struct {
	userInfoStore user_store.IUserInfoStore
	userIdCache   IUserInfoCacheId
	userNameCache IUserInfoCacheName
}

func newUserInfoService() user.IUserInfoService {
	u := &userInfoService{}
	bean.Autowired(&u.userInfoStore)
	bean.Autowired(&u.userIdCache)
	bean.Autowired(&u.userNameCache)
	apinto_module.RegisterEventHandler("login", u.loginHandler)
	return u
}

func decode(v any) (*user_model.UserBase, error) {
	return apinto_module.DecodeFor[user_model.UserBase](v)
}

func (u *userInfoService) save(ctx context.Context, info *user_entry.UserInfo) error {
	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {
		userModel := entryToModule(info)
		err := u.userInfoStore.Save(ctx, info)
		if err != nil {
			return err
		}
		err = u.userNameCache.Set(ctx, userModel.UserName, userModel, time.Hour)
		if err != nil {
			return err
		}
		return u.userIdCache.Set(ctx, userModel.Id, userModel, time.Hour)
	})
}

func (u *userInfoService) UpdateMyPassword(ctx context.Context, userId int, req *user_dto.UpdateMyPasswordReq) error {
	info, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	if common.Md5(req.Old) != info.Password {
		return errors.New("error old password")
	}
	info.Password = common.Md5(req.Password)
	return u.save(ctx, info)
}

func (u *userInfoService) UpdateLastLoginTime(ctx context.Context, userId int, loginTime *time.Time) error {
	info, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	info.LastLoginTime = loginTime
	return u.save(ctx, info)

}

func (u *userInfoService) loginHandler(login string, v any) {
	userBase, err := decode(v)
	if err != nil {
		return
	}
	now := time.Now()
	userEntry := &user_entry.UserInfo{
		Id:            0,
		Sex:           userBase.Sex,
		UserName:      userBase.UserName,
		NoticeUserId:  userBase.NoticeUserId,
		NickName:      userBase.NickName,
		Email:         userBase.Email,
		Phone:         userBase.Phone,
		Avatar:        userBase.Avatar,
		LastLoginTime: &now,
	}
	u.userInfoStore.Save(context.Background(), userEntry)
}

func (u *userInfoService) GetAllUsers(ctx context.Context) ([]*user_model.UserInfo, error) {
	infos, err := u.userInfoStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resList := make([]*user_model.UserInfo, 0, len(infos))

	for _, info := range infos {
		resList = append(resList, &user_model.UserInfo{
			Id:            info.Id,
			Sex:           info.Sex,
			UserName:      info.UserName,
			NoticeUserId:  info.NoticeUserId,
			NickName:      info.NickName,
			Email:         info.Email,
			Phone:         info.Phone,
			Avatar:        info.Avatar,
			LastLoginTime: info.LastLoginTime,
		})
	}

	return resList, nil
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
				}
			}
		} else {
			//获取所有用户并存到缓存
			userList, err := u.userInfoStore.GetAll(ctx)
			if err != nil {
				//只有当数据库报错时才会到这
				//补全返回信息
				tempMaps := make(map[int]*user_model.UserInfo, len(userIds))
				for _, userId := range userIds {
					userModel := &user_model.UserInfo{
						Id:            userId,
						Sex:           0,
						UserName:      "unknown",
						NoticeUserId:  "",
						NickName:      "unknown",
						Email:         "",
						Phone:         "unknown",
						Avatar:        "",
						LastLoginTime: nil,
					}
					tempMaps[userId] = userModel
				}
				return tempMaps, nil
			}
			userSet := make(map[int]struct{}, len(userIds))
			for _, userId := range userIds {
				userSet[userId] = struct{}{}
			}
			tempMaps := make(map[int]*user_model.UserInfo, len(userIds))
			for _, userInfo := range userList {
				userModel := entryToModule(userInfo)
				if _, has := userSet[userInfo.Id]; has {
					tempMaps[userInfo.Id] = userModel
					delete(userSet, userInfo.Id)
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel, time.Hour)
				u.userNameCache.Set(ctx, userModel.UserName, userModel, time.Hour)
			}
			//补全传入的userIds中数据库不存在的数据
			for userID := range userSet {
				userModel := &user_model.UserInfo{
					Id:            userID,
					Sex:           0,
					UserName:      "unknown",
					NoticeUserId:  "",
					NickName:      "unknown",
					Email:         "",
					Phone:         "unknown",
					Avatar:        "",
					LastLoginTime: nil,
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel, time.Hour)
				tempMaps[userID] = userModel
			}
			maps = tempMaps
		}
	}

	return maps, nil
}

func (u *userInfoService) getCache(ctx context.Context, userID int) (*user_model.UserInfo, bool) {
	userModel, err := u.userIdCache.Get(ctx, userID)
	if err != nil {
		return nil, false
	}
	return userModel, true
}

func (u *userInfoService) getCacheByName(ctx context.Context, userName string) (*user_model.UserInfo, bool) {
	userModel, err := u.userNameCache.Get(ctx, userName)
	if err != nil {
		return nil, false
	}
	return userModel, true
}

// GetUserInfo 获取不到用户信息记录错误即可，不必返回error
func (u *userInfoService) GetUserInfo(ctx context.Context, userID int) (*user_model.UserInfo, error) {
	userModel, err := u.userIdCache.Get(ctx, userID)
	if err == nil {
		return userModel, nil
	}
	//var userModel *user_model.UserInfo
	userInfo, err := u.userInfoStore.Get(ctx, userID)
	if err != nil {
		userModel = &user_model.UserInfo{
			Id:            userID,
			Sex:           0,
			UserName:      "unknown",
			NoticeUserId:  "",
			NickName:      "unknown",
			Email:         "",
			Phone:         "unknown",
			Avatar:        "",
			LastLoginTime: nil,
		}
	} else {
		userModel = entryToModule(userInfo)
		u.userNameCache.Set(ctx, userModel.UserName, userModel, time.Hour)
	}
	u.userIdCache.Set(ctx, userID, userModel, time.Hour)
	return userModel, nil
}

func (u *userInfoService) CheckPassword(ctx context.Context, name string, password string) (int, bool) {
	//var userModel *user_model.UserInfo
	info, err := u.userInfoStore.GetByUserName(ctx, name)
	if err != nil {
		return 0, false
	}
	return info.Id, common.Md5(password) != info.Password
}

func (u *userInfoService) GetUserInfoByName(ctx context.Context, userName string) (*user_model.UserInfo, error) {

	userModel, err := u.userNameCache.Get(ctx, userName)
	if err == nil {
		return userModel, nil
	}
	//var userModel *user_model.UserInfo
	userInfo, err := u.userInfoStore.GetByUserName(ctx, userName)
	if err != nil {
		userModel = &user_model.UserInfo{
			Id:            0,
			Sex:           0,
			UserName:      userName,
			NoticeUserId:  "",
			NickName:      "unknown",
			Email:         "",
			Phone:         "unknown",
			Avatar:        "",
			LastLoginTime: nil,
		}
	} else {
		userModel = entryToModule(userInfo)
		u.userIdCache.Set(ctx, userModel.Id, userModel, time.Hour)
	}
	u.userNameCache.Set(ctx, userName, userModel, time.Hour)
	return userModel, nil
}

func (u *userInfoService) GetUserInfoByNames(ctx context.Context, userNames ...string) (map[string]*user_model.UserInfo, error) {
	maps := make(map[string]*user_model.UserInfo)
	need := make([]string, 0, len(userNames))

	for _, userName := range userNames {
		userInfo, ok := u.getCacheByName(ctx, userName)
		if ok {
			maps[userInfo.UserName] = userInfo
		} else {
			need = append(need, userName)
		}
	}
	if len(need) > 0 {
		if len(need) < 5 {
			for _, userId := range need {
				userInfo, err := u.GetUserInfoByName(ctx, userId)
				if err == nil {
					maps[userInfo.UserName] = userInfo
				}
			}
		} else {
			//获取所有用户并存到缓存
			userList, err := u.userInfoStore.GetAll(ctx)
			if err != nil {
				//只有当数据库报错时才会到这
				//补全返回信息
				tempMaps := make(map[string]*user_model.UserInfo, len(userNames))
				for _, userName := range userNames {
					userModel := &user_model.UserInfo{
						Id:            0,
						Sex:           0,
						UserName:      userName,
						NoticeUserId:  "",
						NickName:      "unknown",
						Email:         "",
						Phone:         "unknown",
						Avatar:        "",
						LastLoginTime: nil,
					}
					tempMaps[userName] = userModel
				}
				return tempMaps, nil
			}
			userSet := make(map[string]struct{}, len(userNames))
			for _, userName := range userNames {
				userSet[userName] = struct{}{}
			}
			tempMaps := make(map[string]*user_model.UserInfo, len(userNames))
			for _, userInfo := range userList {
				userModel := entryToModule(userInfo)
				if _, has := userSet[userInfo.UserName]; has {
					tempMaps[userInfo.UserName] = userModel
					delete(userSet, userInfo.UserName)
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel, time.Hour)
				u.userNameCache.Set(ctx, userModel.UserName, userModel, time.Hour)
			}
			//补全传入的userIds中数据库不存在的数据
			for userName := range userSet {
				userModel := &user_model.UserInfo{
					Id:            0,
					Sex:           0,
					UserName:      userName,
					NoticeUserId:  "",
					NickName:      "unknown",
					Email:         "",
					Phone:         "unknown",
					Avatar:        "",
					LastLoginTime: nil,
				}
				u.userNameCache.Set(ctx, userModel.UserName, userModel, time.Hour)
				tempMaps[userName] = userModel
			}
			maps = tempMaps
		}
	}

	return maps, nil
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

func (u *userInfoService) UpdateMyProfile(ctx context.Context, userId int, req *user_dto.UpdateMyProfileReq) error {
	info, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	info.NickName = req.NickName
	info.Email = req.Email
	info.NoticeUserId = req.NoticeUserId
	return u.save(ctx, info)

}
