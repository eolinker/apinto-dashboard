package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/user"
	user_entry "github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	user_store "github.com/eolinker/apinto-dashboard/modules/user/user-store"
	"github.com/eolinker/eosc/common/bean"
)

const AdminName = "admin"
const defaultPwd = "12345678"

type userInfoService struct {
	userInfoStore user_store.IUserInfoStore
	userIdCache   IUserInfoCacheId
	userNameCache IUserInfoCacheName
}

func (u *userInfoService) SetUserName(ctx context.Context, iterate ...user.UserInfoIterate) error {
	if len(iterate) == 0 {
		return nil
	}

	ids := make([]int, 0, len(iterate))
	cache := make(map[int]string)
	for _, i := range iterate {
		id := i.UserId()
		c, err := u.userIdCache.Get(ctx, id)
		if err != nil || c == nil {
			ids = append(ids, id)
		} else {
			if c.NickName == "" {
				cache[id] = c.UserName
			} else {
				cache[id] = c.NickName
			}
		}
	}
	if len(ids) > 10 {
		all, err := u.userInfoStore.GetAll(ctx)
		if err != nil {
			return err
		}
		cache = common.SliceToMapO(all, func(t *user_entry.UserInfo) (int, string) {
			n := t.NickName
			if n == "" {
				n = t.UserName
			}
			if _, has := cache[t.Id]; !has {
				u.userIdCache.Set(ctx, t.Id, user_model.CreateUserInfo(t))
			}
			return t.Id, n
		})
	} else {
		for _, id := range ids {
			en, err := u.userInfoStore.Get(ctx, id)
			if err != nil {
				continue
			}
			if en.NickName == "" {
				cache[id] = en.UserName
			} else {
				cache[id] = en.NickName
			}
			u.userIdCache.Set(ctx, id, user_model.CreateUserInfo(en))
		}
	}
	for _, i := range iterate {
		i.Set(cache[i.UserId()])
	}

	return nil
}

func (u *userInfoService) SaveUserInfo(ctx context.Context, info *user_model.UserInfo) error {
	if info.Id == 0 {
		return nil
	}
	entity := &user_entry.UserInfo{
		Id:           info.Id,
		Sex:          info.Sex,
		UserName:     info.UserName,
		NoticeUserId: info.NoticeUserId,
		NickName:     info.NickName,
		Email:        info.Email,
		Phone:        info.Phone,
		Avatar:       info.Avatar,
	}
	_, err := u.userInfoStore.DeleteWhere(ctx, map[string]interface{}{
		"username": info.UserName,
	})
	if err != nil {
		return err
	}
	err = u.userInfoStore.Save(ctx, entity)
	if err != nil {
		return err
	}
	u.userNameCache.Set(ctx, info.UserName, info)
	u.userIdCache.Set(ctx, info.Id, info)
	return nil

}

func newUserInfoService() user.IUserInfoService {
	u := &userInfoService{}
	bean.Autowired(&u.userInfoStore)
	bean.Autowired(&u.userIdCache)
	bean.Autowired(&u.userNameCache)
	//apinto_module.RegisterEventHandler("login", u.loginHandler)
	//apinto_module.RegisterEventHandler("user-create", u.userUpdate)
	//apinto_module.RegisterEventHandler("user-update", u.userUpdate)
	//apinto_module.RegisterEventHandler("user-delete", u.userDelete)
	return u
}

func (u *userInfoService) save(ctx context.Context, info *user_entry.UserInfo) error {
	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {

		err := u.userInfoStore.Save(ctx, info)
		if err != nil {
			return err
		}
		userModel := user_model.CreateUserInfo(info)
		err = u.userNameCache.Set(ctx, userModel.UserName, userModel)
		if err != nil {
			return err
		}
		return u.userIdCache.Set(ctx, userModel.Id, userModel)
	})
}

//
//func (u *userInfoService) loginHandler(login string, v any) {
//	userBase, err := decode(v)
//	if err != nil {
//		return
//	}
//	now := time.Now()
//	userEntry := &user_entry.UserInfo{
//		Id:            0,
//		Sex:           userBase.Sex,
//		UserName:      userBase.UserName,
//		NoticeUserId:  userBase.NoticeUserId,
//		NickName:      userBase.NickName,
//		Email:         userBase.Email,
//		Phone:         userBase.Phone,
//		Avatar:        userBase.Avatar,
// 	}
//	u.userInfoStore.Save(context.Background(), userEntry)
//}
//func (u *userInfoService) userUpdate(event string, v any) {
//	userBase, err := decode(v)
//	if err != nil {
//		return
//	}
//	userEntry, err := u.userInfoStore.GetByUserName(context.Background(), userBase.UserName)
//	if err != nil {
//		now := time.Now()
//		userEntry = &user_entry.UserInfo{
//			Id:            0,
//			Sex:           userBase.Sex,
//			UserName:      userBase.UserName,
//			NoticeUserId:  userBase.NoticeUserId,
//			NickName:      userBase.NickName,
//			Email:         userBase.Email,
//			Phone:         userBase.Phone,
//			Avatar:        userBase.Avatar,
//			LastLoginTime: &now,
//		}
//		u.userInfoStore.Save(context.Background(), userEntry)
//		return
//	}
//	userEntry.Sex = userBase.Sex
//	userEntry.Email = userBase.Email
//	userEntry.NickName = userBase.NickName
//	userEntry.NoticeUserId = userBase.NoticeUserId
//	userEntry.Avatar = userBase.Avatar
//	userEntry.Phone = userBase.Phone
//	u.userInfoStore.Update(context.Background(), userEntry)
//	u.userNameCache.Set(context.Background(), userEntry.UserName, user_model.CreateUserInfo(userEntry))
//
//}
//func (u *userInfoService) userDelete(e string, v any) {
//	return
//}

func (u *userInfoService) GetAllUsers(ctx context.Context) ([]*user_model.UserInfo, error) {
	infos, err := u.userInfoStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resList := make([]*user_model.UserInfo, 0, len(infos))

	for _, info := range infos {
		resList = append(resList, user_model.CreateUserInfo(info))
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
						Id:           userId,
						Sex:          0,
						UserName:     "unknown",
						NoticeUserId: "",
						NickName:     "unknown",
						Email:        "",
						Phone:        "unknown",
						Avatar:       "",
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
				userModel := user_model.CreateUserInfo(userInfo)
				if _, has := userSet[userInfo.Id]; has {
					tempMaps[userInfo.Id] = userModel
					delete(userSet, userInfo.Id)
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel)
				u.userNameCache.Set(ctx, userModel.UserName, userModel)
			}
			//补全传入的userIds中数据库不存在的数据
			for userID := range userSet {
				userModel := &user_model.UserInfo{
					Id:           userID,
					Sex:          0,
					UserName:     "unknown",
					NoticeUserId: "",
					NickName:     "unknown",
					Email:        "",
					Phone:        "unknown",
					Avatar:       "",
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel)
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
			Id:           userID,
			Sex:          0,
			UserName:     "unknown",
			NoticeUserId: "",
			NickName:     "unknown",
			Email:        "",
			Phone:        "unknown",
			Avatar:       "",
		}
	} else {
		userModel = user_model.CreateUserInfo(userInfo)
		u.userNameCache.Set(ctx, userModel.UserName, userModel)
	}
	u.userIdCache.Set(ctx, userID, userModel)
	return userModel, nil
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
			Id:           0,
			Sex:          0,
			UserName:     userName,
			NoticeUserId: "",
			NickName:     "unknown",
			Email:        "",
			Phone:        "unknown",
			Avatar:       "",
		}
	} else {
		userModel = user_model.CreateUserInfo(userInfo)
		u.userIdCache.Set(ctx, userModel.Id, userModel)
	}
	u.userNameCache.Set(ctx, userName, userModel)
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
						Id:           0,
						Sex:          0,
						UserName:     userName,
						NoticeUserId: "",
						NickName:     "unknown",
						Email:        "",
						Phone:        "unknown",
						Avatar:       "",
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
				userModel := user_model.CreateUserInfo(userInfo)
				if _, has := userSet[userInfo.UserName]; has {
					tempMaps[userInfo.UserName] = userModel
					delete(userSet, userInfo.UserName)
				}
				u.userIdCache.Set(ctx, userModel.Id, userModel)
				u.userNameCache.Set(ctx, userModel.UserName, userModel)
			}
			//补全传入的userIds中数据库不存在的数据
			for userName := range userSet {
				userModel := &user_model.UserInfo{
					Id:           0,
					Sex:          0,
					UserName:     userName,
					NoticeUserId: "",
					NickName:     "unknown",
					Email:        "",
					Phone:        "unknown",
					Avatar:       "",
				}
				u.userNameCache.Set(ctx, userModel.UserName, userModel)
				tempMaps[userName] = userModel
			}
			maps = tempMaps
		}
	}

	return maps, nil
}

//func (u *userInfoService) UpdateMyProfile(ctx context.Context, userId int, req *user_dto.UpdateMyProfileReq) error {
//	info, err := u.userInfoStore.Get(ctx, userId)
//	if err != nil {
//		return err
//	}
//
//	info.NickName = req.NickName
//	info.Email = req.Email
//	info.NoticeUserId = req.NoticeUserId
//	return u.save(ctx, info)
//
//}
