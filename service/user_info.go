package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/apinto-dashboard/user_center/client"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

const AdminName = "admin"
const defaultPwd = "12345678"

var iv = []byte("1e42=7838a1vfc6n")

type IUserInfoService interface {
	GetUserInfo(ctx context.Context, userId int) (*model.UserInfo, error)
	UpdateMyProfile(ctx context.Context, userId int, req *dto.UpdateMyProfileReq) error
	UpdateMyPassword(ctx context.Context, userId int, req *dto.UpdateMyPasswordReq) error
	UpdateLoginTime(ctx context.Context, userId int)
	GetUserInfoList(ctx context.Context, roleId, keyWord string) ([]*model.UserInfo, error)
	GetUserInfoAll(ctx context.Context) ([]*model.UserInfo, error)
	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*model.UserInfo, error)
	CreateUser(ctx context.Context, operator int, userInfo *dto.SaveUserReq) error
	CheckUser(ctx context.Context, userId int) error
	PatchUser(ctx context.Context, operator, userId int, req *dto.PatchUserReq) error
	UpdateUser(ctx context.Context, operator, userId int, userInfo *dto.SaveUserReq) error
	DelUser(ctx context.Context, operator int, userIds []int) error
	ResetUserPwd(ctx context.Context, operator, userId int, pwd string) error

	GetAccessInfo(ctx context.Context, userId int) (map[access.Access]struct{}, error)
	GetRoleList(ctx context.Context, userID int) ([]*model.RoleListItem, int, error)
	GetRoleInfo(ctx context.Context, roleID string) (*model.RoleInfo, error)
	GetRoleOptions(ctx context.Context) ([]*model.RoleOptionItem, error)
	CreateRole(ctx context.Context, userID int, input *dto.ProxyRoleInfo) error
	UpdateRole(ctx context.Context, userID int, roleUUID string, input *dto.ProxyRoleInfo) error
	DeleteRole(ctx context.Context, userID int, roleUUID string) error
	RoleBatchUpdate(ctx context.Context, userIds []int, roleUUID string) error
	RoleBatchRemove(ctx context.Context, userIds []int, roleUUID string) error

	GetUserRoleIds(ctx context.Context, userID int) ([]string, error)
	GetRoleAccessIds(ctx context.Context, roleUUID string) ([]string, error)
	CreateAdmin() error
	CleanAdminCache()
}

type userInfoService struct {
	userInfoStore      store.IUserInfoStore
	userRoleStore      store.IUserRoleStore
	roleStore          store.IRoleStore
	roleAccessStore    store.IRoleAccessStore
	roleAccessLogStore store.IRoleAccessLogStore

	userInfoCache    cache.IUserInfoCache
	roleAccessCache  cache.IRoleAccessCache
	userCenterClient client.IUserCenterClient
}

func (u *userInfoService) CleanAdminCache() {
	for key, _ := range access.GetBuildInRoleMap() {
		_ = u.roleAccessCache.Delete(context.TODO(), u.roleAccessCache.Key(key))
	}
}

func newUserInfoService() IUserInfoService {
	u := &userInfoService{}

	bean.Autowired(&u.userInfoStore)
	bean.Autowired(&u.userInfoCache)
	bean.Autowired(&u.roleAccessCache)
	bean.Autowired(&u.userCenterClient)
	bean.Autowired(&u.userRoleStore)
	bean.Autowired(&u.roleStore)
	bean.Autowired(&u.roleAccessStore)
	bean.Autowired(&u.roleAccessLogStore)
	return u
}

// CreateAdmin 创建超管
func (u *userInfoService) CreateAdmin() error {
	//如果已经有超管，就不需要创建了
	ctx := context.Background()
	currentUserInfo, err := u.userInfoStore.GetByUserName(ctx, AdminName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound || currentUserInfo.IsDelete {
		//创建超管
		return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {

			password, err := bcrypt.GenerateFromPassword([]byte(common.Md5(defaultPwd)), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			createUser := &client.CreateUserReq{
				UserNickName: AdminName,
				UserName:     AdminName,
				Password:     string(password),
				Status:       1,
			}
			userId, err := u.userCenterClient.CreateUser(createUser)
			if err != nil {
				return err
			}

			t := time.Now()

			adminRole := access.GetAdminRole()
			if adminRole == nil {
				return errors.New("获取超管权限失败")
			}

			userInfo := &entry.UserInfo{
				Id:         *userId,
				UserName:   AdminName,
				NickName:   AdminName,
				RoleIds:    adminRole.Uuid,
				Status:     2,
				FlushTime:  t,
				CreateTime: t,
				UpdateTime: t,
			}

			if err = u.userInfoStore.Save(txCtx, userInfo); err != nil {
				return err
			}

			userRole := &entry.UserRole{
				UserID:     userInfo.Id,
				RoleID:     adminRole.ID,
				Module:     "/",
				CreateTime: t,
			}

			if err = u.userRoleStore.Save(txCtx, userRole); err != nil {
				return err
			}

			return nil
		})
	}
	return nil
}

func (u *userInfoService) UpdateMyPassword(ctx context.Context, userId int, req *dto.UpdateMyPasswordReq) error {
	userInfoReq := &client.UserInfoReq{
		Id: strconv.Itoa(userId),
	}
	info, err := u.userCenterClient.UserInfo(userInfoReq)
	if err != nil {
		return err
	}

	key := common.Md5(info.UserName)

	decodeOldPwd, err := common.Base64Decode(req.Old)
	if err != nil {
		return err
	}

	oldPwdDecrypter := common.CBCDecrypter(decodeOldPwd, []byte(key), iv)

	decodeNewPwd, err := common.Base64Decode(req.Password)
	if err != nil {
		return err
	}

	newPwdDecrypter := common.CBCDecrypter(decodeNewPwd, []byte(key), iv)

	updateUserPwdReq := &client.UpdateUserPwdReq{
		UserId:      userId,
		AccountName: info.UserName,
		OldPassword: string(oldPwdDecrypter),
		Password:    string(newPwdDecrypter),
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(info.Id),
		Name: info.UserName,
	})

	return u.userCenterClient.UpdateUserPwd(updateUserPwdReq)
}

func (u *userInfoService) UpdateLoginTime(ctx context.Context, userId int) {

	err := u.userInfoStore.UpdateLoginTime(ctx, userId, time.Now())
	if err != nil {
		log.Errorf("userInfoService-updateLoginTime userId=%d err=%s", userId, err.Error())
	}
}

func (u *userInfoService) UpdateMyProfile(ctx context.Context, userId int, req *dto.UpdateMyProfileReq) error {
	info, err := u.GetUserInfo(ctx, userId)
	if err != nil {
		return err
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(info.Id),
		Name: info.UserName,
	})

	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {
		info.NickName = req.NickName
		info.Email = req.Email
		info.Remark = req.Desc
		info.NoticeUserId = req.NoticeUserId

		if err = u.userInfoStore.Save(txCtx, info.UserInfo); err != nil {
			return err
		}

		updateUserReq := &client.UpdateUserReq{
			AccountName:  info.UserName,
			UserNickName: info.NickName,
			Email:        info.Email,
			Remark:       info.Remark,
		}

		if err = u.userCenterClient.UpdateUser(updateUserReq); err != nil {
			return err
		}

		_ = u.userInfoCache.Set(ctx, u.userInfoCache.Key(userId), info.UserInfo, time.Minute*30)

		return nil
	})

}

func (u *userInfoService) PatchUser(ctx context.Context, operator, userId int, req *dto.PatchUserReq) error {
	if len(req.Role) == 0 && req.Status == 0 {
		return errors.New("参数错误")
	}

	userInfo, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	if req.Status > 0 {

		err = u.userInfoStore.UpdateStatus(ctx, userId, req.Status)
		if err != nil {
			return err
		}
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(userInfo.Id),
		Name: userInfo.UserName,
	})

	if len(req.Role) > 0 && req.Role[0] != "" {

		roles, err := u.roleStore.GetByUUIDS(ctx, req.Role)
		if err != nil {
			return err
		}

		return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {

			t := time.Now()

			userInfo.UpdateTime = t
			userInfo.RoleIds = strings.Join(req.Role, ",")

			if err = u.userInfoStore.Save(txCtx, userInfo); err != nil {
				return err
			}

			if err = u.userRoleStore.DelByUserId(txCtx, userId); err != nil {
				return err
			}

			for _, role := range roles {
				userRole := &entry.UserRole{
					UserID:     userId,
					RoleID:     role.ID,
					CreateTime: t,
				}
				if err = u.userRoleStore.Save(txCtx, userRole); err != nil {
					return err
				}
			}

			_ = u.userInfoCache.Set(ctx, u.userInfoCache.Key(userId), userInfo, time.Minute*30)

			return nil
		})

	}

	return nil
}

func (u *userInfoService) CheckUser(ctx context.Context, userId int) error {
	//检测用户是否可登录
	userInfo, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}
	if userInfo.Status == 1 {
		return errors.New("该账号已被禁用")
	}
	if userInfo.IsDelete {
		return errors.New("该账号已被删除")
	}

	return nil
}

func (u *userInfoService) ResetUserPwd(ctx context.Context, operator, userId int, pwd string) error {

	userInfo, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(common.Md5(pwd)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req := &client.ResetUserPwd{
		UserId:   userId,
		Password: string(password),
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(userId),
		Name: userInfo.UserName,
	})

	return u.userCenterClient.ResetUserPwd(req)
}

func (u *userInfoService) GetUserInfoAll(ctx context.Context) ([]*model.UserInfo, error) {

	infos, err := u.userInfoStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	resList := make([]*model.UserInfo, 0, len(infos))

	for _, info := range infos {
		resList = append(resList, &model.UserInfo{
			UserInfo: info,
		})
	}

	return resList, nil
}

func (u *userInfoService) GetUserInfoList(ctx context.Context, roleUUID, keyword string) ([]*model.UserInfo, error) {

	roleId := 0

	isBuildInRole := access.GetBuildInRole(roleUUID)

	if roleUUID != "" && isBuildInRole == nil {
		role, err := u.roleStore.GetByUUID(ctx, roleUUID)
		if err != nil {
			return nil, err
		}
		roleId = role.ID
	}
	if isBuildInRole != nil {
		roleId = isBuildInRole.ID
	}

	infos, err := u.userInfoStore.GetByRoleByKeyword(ctx, roleId, keyword)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(infos, func(t *entry.UserInfo) int {
		return t.Operator
	})

	userInfoMaps, _ := u.GetUserInfoMaps(ctx, userIds...)

	resList := make([]*model.UserInfo, 0, len(infos))

	for _, info := range infos {

		operatorName := ""
		if userInfo, ok := userInfoMaps[info.Operator]; ok {
			operatorName = userInfo.NickName
		}

		resList = append(resList, &model.UserInfo{
			UserInfo:      info,
			OperateEnable: info.UserName == AdminName,
			Operator:      operatorName,
		})
	}

	return resList, nil
}

func (u *userInfoService) DelUser(ctx context.Context, operator int, userIds []int) error {

	infoMaps, _ := u.GetUserInfoMaps(ctx, userIds...)
	for _, info := range infoMaps {
		if info.UserName == AdminName {
			return errors.New("超管不可删除")
		}
	}

	if len(infoMaps) != len(userIds) {
		return errors.New("修改过程中用户数据出现了不一致，请稍后再试")
	}

	uuids := make([]string, 0)
	names := make([]string, 0)
	for _, info := range infoMaps {
		uuids = append(uuids, strconv.Itoa(info.Id))
		names = append(names, info.UserName)
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{Uuid: strings.Join(uuids, ","), Name: strings.Join(names, ",")})

	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {

		if err := u.userInfoStore.SoftDelete(txCtx, operator, userIds...); err != nil {
			return err
		}

		if err := u.userCenterClient.DelUser(userIds...); err != nil {
			return err
		}

		for _, id := range userIds {

			if err := u.userRoleStore.DelByUserId(txCtx, id); err != nil {
				return err
			}

			u.userInfoCache.Delete(ctx, u.userInfoCache.Key(id))
		}
		return nil
	})

}

func (u *userInfoService) GetUserInfoMaps(ctx context.Context, userIds ...int) (map[int]*model.UserInfo, error) {

	maps := make(map[int]*model.UserInfo)
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
func (u *userInfoService) GetUserInfo(ctx context.Context, userID int) (*model.UserInfo, error) {

	key := u.userInfoCache.Key(userID)

	//读缓存
	userInfo, err := u.userInfoCache.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if userInfo != nil {
		return &model.UserInfo{UserInfo: userInfo}, err
	}

	//读数据库
	userInfo, err = u.userInfoStore.Get(ctx, userID)
	if err == gorm.ErrRecordNotFound {
		_ = u.userInfoCache.Set(ctx, key, &entry.UserInfo{}, time.Minute*30)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	//重新放到缓存
	_ = u.userInfoCache.Set(ctx, key, userInfo, time.Minute*30)

	return &model.UserInfo{UserInfo: userInfo}, nil
}

func (u *userInfoService) GetUserInfoByUUID(ctx context.Context, userID string) (*model.UserInfo, error) {

	userInfo, err := u.userInfoStore.GetByUUID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserInfo{UserInfo: userInfo}, nil
}

func (u *userInfoService) CreateUser(ctx context.Context, operator int, userInfo *dto.SaveUserReq) error {

	roles, err := u.roleStore.GetByUUIDS(ctx, userInfo.RoleIds)
	if err != nil {
		return err
	}
	_, err = u.userInfoStore.GetByUserName(ctx, userInfo.UserName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err != gorm.ErrRecordNotFound {
		return errors.New("用户已存在")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(common.Md5(defaultPwd)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userCenterSaveUserReq := &client.CreateUserReq{
		OperatorId:   operator,
		UserNickName: userInfo.NickName,
		UserName:     userInfo.UserName,
		Sex:          userInfo.Sex,
		Avatar:       userInfo.Avatar,
		Email:        userInfo.Email,
		Phone:        userInfo.Phone,
		Password:     string(password),
		Status:       1,
		Remark:       userInfo.Desc,
	}

	userId, err := u.userCenterClient.CreateUser(userCenterSaveUserReq)
	if err != nil {
		return err
	}

	t := time.Now()
	entryUserInfo := &entry.UserInfo{
		Id:           *userId,
		Sex:          userInfo.Sex,
		UserName:     userInfo.UserName,
		NoticeUserId: userInfo.NoticeUserId,
		NickName:     userInfo.NickName,
		Email:        userInfo.Email,
		Phone:        userInfo.Phone,
		Avatar:       userInfo.Avatar,
		Remark:       userInfo.Desc,
		RoleIds:      strings.Join(userInfo.RoleIds, ","),
		Status:       2,
		Operator:     operator,
		FlushTime:    t,
		CreateTime:   t,
		UpdateTime:   t,
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(*userId),
		Name: userInfo.UserName,
	})

	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = u.userInfoStore.Save(txCtx, entryUserInfo); err != nil {
			return err
		}

		userRoles := make([]*entry.UserRole, 0)

		for _, role := range roles {
			userRole := &entry.UserRole{
				UserID:     entryUserInfo.Id,
				RoleID:     role.ID,
				Module:     "/",
				CreateTime: t,
			}
			userRoles = append(userRoles, userRole)
		}

		if len(userRoles) > 0 {
			if err = u.userRoleStore.Insert(txCtx, userRoles...); err != nil {
				return err
			}
		}

		return u.userInfoCache.Set(ctx, u.userInfoCache.Key(*userId), entryUserInfo, time.Minute*30)
	})

}

func (u *userInfoService) UpdateUser(ctx context.Context, operator, userId int, userInfo *dto.SaveUserReq) error {

	roles, err := u.roleStore.GetByUUIDS(ctx, userInfo.RoleIds)
	if err != nil {
		return err
	}

	updateUserReq := &client.UpdateUserReq{
		AccountName:  userInfo.UserName,
		Avatar:       userInfo.Avatar,
		UserNickName: userInfo.NickName,
		Sex:          userInfo.Sex,
		Email:        userInfo.Email,
		UserName:     userInfo.UserName,
		Phone:        userInfo.Phone,
		Remark:       userInfo.Desc,
	}

	if err = u.userCenterClient.UpdateUser(updateUserReq); err != nil {
		return err
	}

	info, err := u.userInfoStore.Get(ctx, userId)
	if err != nil {
		return err
	}

	t := time.Now()
	newUserInfo := &entry.UserInfo{
		Id:           info.Id,
		Sex:          userInfo.Sex,
		UserName:     userInfo.UserName,
		NoticeUserId: userInfo.NoticeUserId,
		NickName:     userInfo.NickName,
		Email:        userInfo.Email,
		Phone:        userInfo.Phone,
		Avatar:       userInfo.Avatar,
		Remark:       userInfo.Desc,
		RoleIds:      strings.Join(userInfo.RoleIds, ","),
		Status:       info.Status,
		Operator:     operator,
		FlushTime:    info.FlushTime,
		CreateTime:   t,
		UpdateTime:   t,
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: strconv.Itoa(info.Id),
		Name: userInfo.UserName,
	})

	//超管不能更改角色
	if userInfo.UserName == AdminName {
		newUserInfo.RoleIds = info.RoleIds
	}

	return u.userInfoStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = u.userInfoStore.Save(txCtx, newUserInfo); err != nil {
			return err
		}

		if userInfo.UserName != AdminName {
			//先删除用户的角色 然后添加新的角色
			if err = u.userRoleStore.DelByUserId(txCtx, newUserInfo.Id); err != nil {
				return err
			}

			userRoles := make([]*entry.UserRole, 0)

			for _, role := range roles {
				userRole := &entry.UserRole{
					UserID:     newUserInfo.Id,
					RoleID:     role.ID,
					Module:     "/",
					CreateTime: t,
				}
				userRoles = append(userRoles, userRole)
			}

			if len(userRoles) > 0 {
				if err = u.userRoleStore.Insert(txCtx, userRoles...); err != nil {
					return err
				}

			}

		}

		//保存到缓存
		return u.userInfoCache.Set(ctx, u.userInfoCache.Key(info.Id), newUserInfo, time.Minute*30)
	})

}

func (u *userInfoService) GetAccessInfo(ctx context.Context, userId int) (map[access.Access]struct{}, error) {
	accessSet := make(map[access.Access]struct{})

	roleIds, err := u.GetUserRoleIds(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, role := range roleIds {
		//若缓存没有，则查表，并存入缓存
		accessIds, err := u.GetRoleAccessIds(ctx, role)
		if err != nil {
			return nil, err
		}
		for _, accessId := range accessIds {
			if accessId == "" || len(accessId) == 0 {
				continue
			}
			id, _ := strconv.Atoi(accessId)
			accessSet[access.Access(id)] = struct{}{}
		}
	}
	return accessSet, nil
}

func (u *userInfoService) GetRoleList(ctx context.Context, userID int) ([]*model.RoleListItem, int, error) {
	//获取内置角色
	buildInRoles := access.GetBuildInRoleList()
	//获取数据库表中的自定义角色
	roles, err := u.roleStore.GetAllRole(ctx)
	if err != nil {
		return nil, 0, err
	}

	roleList := make([]*model.RoleListItem, 0, len(buildInRoles)+len(roles))

	for _, role := range buildInRoles {
		userNum, err := u.userRoleStore.GetRoleQuotedCount(ctx, role.ID)
		if err != nil {
			return nil, 0, err
		}

		item := &model.RoleListItem{
			ID:             role.Uuid,
			Title:          role.Title,
			UserNum:        userNum,
			OperateDisable: !role.IsAddable,
			Type:           1, //内置角色
		}
		roleList = append(roleList, item)
	}

	for _, role := range roles {
		//从用户角色表获取该角色有多少个用户
		userNum, err := u.userRoleStore.GetRoleQuotedCount(ctx, role.ID)
		if err != nil {
			return nil, 0, err
		}

		item := &model.RoleListItem{
			ID:             role.Uuid,
			Title:          role.Title,
			UserNum:        userNum,
			OperateDisable: false,
			Type:           2, //自定义角色
		}

		roleList = append(roleList, item)
	}
	//获取所有用户的数量
	total, err := u.userInfoStore.GetUserAmount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return roleList, total, nil
}

func (u *userInfoService) GetRoleInfo(ctx context.Context, roleUUID string) (*model.RoleInfo, error) {
	if access.IsBuildInRole(roleUUID) {
		role := access.GetBuildInRole(roleUUID)
		return &model.RoleInfo{
			Title:  role.Title,
			Desc:   role.Desc,
			Access: role.Access,
		}, nil
	}

	role, err := u.roleStore.GetByUUID(ctx, roleUUID)
	if err != nil {
		return nil, err
	}

	roleAccess, err := u.roleAccessStore.GetByRoleID(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	accessIds := strings.Split(roleAccess.AccessIDs, ",")

	accessList := make([]string, 0, len(accessIds))
	for _, strID := range accessIds {
		accessID, err := strconv.Atoi(strID)
		if err != nil {
			continue
		}
		accessList = append(accessList, access.Access(accessID).Key())
	}

	roleInfo := &model.RoleInfo{
		Title:  role.Title,
		Desc:   role.Desc,
		Access: accessList,
	}

	return roleInfo, nil
}

func (u *userInfoService) GetRoleOptions(ctx context.Context) ([]*model.RoleOptionItem, error) {
	//获取内置角色
	buildInRoles := access.GetBuildInRoleList()
	//获取数据库表中的自定义角色
	roles, err := u.roleStore.GetAllRole(ctx)
	if err != nil {
		return nil, err
	}

	roleList := make([]*model.RoleOptionItem, 0, len(buildInRoles)+len(roles))

	for _, role := range buildInRoles {
		item := &model.RoleOptionItem{
			ID:             role.Uuid,
			Title:          role.Title,
			OperateDisable: !role.IsAddable,
		}
		roleList = append(roleList, item)
	}

	for _, role := range roles {
		item := &model.RoleOptionItem{
			ID:             role.Uuid,
			Title:          role.Title,
			OperateDisable: false,
		}

		roleList = append(roleList, item)
	}
	return roleList, nil
}

func (u *userInfoService) CreateRole(ctx context.Context, operator int, input *dto.ProxyRoleInfo) error {
	//判断title有没有重复
	role, err := u.roleStore.GetByTitle(ctx, input.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	//若新建的title与自定义角色或内置角色重名
	if role != nil || access.IsBuildInRoleTitle(input.Title) {
		return fmt.Errorf("title %s is reduplicatd. ", input.Title)
	}

	roleUuid := uuid.New()

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: roleUuid,
		Name: input.Title,
	})

	return u.roleStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		roleInfo := &entry.Role{
			Title:      input.Title,
			Uuid:       roleUuid,
			Desc:       input.Desc,
			Operator:   operator,
			Type:       1,
			Module:     "/", //TODO
			CreateTime: t,
			UpdateTime: t,
		}
		if err = u.roleStore.Save(txCtx, roleInfo); err != nil {
			return err
		}

		accessIDs := make([]string, 0, len(input.Access))
		for _, key := range input.Access {
			id, err := access.Parse(key)
			if err != nil {
				return fmt.Errorf("access %s doesn't exist. ", key)
			}
			accessIDs = append(accessIDs, strconv.Itoa(int(id)))
		}

		roleAccess := &entry.RoleAccess{
			RoleID:     roleInfo.ID,
			AccessIDs:  strings.Join(accessIDs, ","),
			CreateTime: t,
		}

		err = u.roleAccessStore.Save(txCtx, roleAccess)
		if err != nil {
			return err
		}
		//角色权限变更日志
		return u.roleAccessLogStore.HistoryAdd(txCtx, 0, roleInfo.ID, &entry.AccessListLog{
			AccessIds: accessIDs,
			Module:    roleInfo.Module,
		}, operator)
	})
}

func (u *userInfoService) UpdateRole(ctx context.Context, userID int, roleUUID string, input *dto.ProxyRoleInfo) error {
	//内置角色不可以修改
	if access.IsBuildInRole(roleUUID) {
		return errors.New("build-in role can't be updated. ")
	}

	//判断角色uuid存不存在
	role, err := u.roleStore.GetByUUID(ctx, roleUUID)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("role_id is invalid. ")
		}
		return err
	}

	//判断title有没有重复
	oldTitleRole, err := u.roleStore.GetByTitle(ctx, input.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	//若旧title已有使用且不是角色本身的title 或者是内置角色的名字
	if (oldTitleRole != nil && oldTitleRole.ID != role.ID) || access.IsBuildInRoleTitle(input.Title) {
		return fmt.Errorf("title %s is reduplicatd. ", input.Title)
	}

	t := time.Now()

	role.Title = input.Title
	role.Desc = input.Desc
	role.UpdateTime = t

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: roleUUID,
		Name: input.Title,
	})

	return u.roleStore.Transaction(ctx, func(txCtx context.Context) error {

		if _, err = u.roleStore.Update(txCtx, role); err != nil {
			return err
		}

		accessIDs := make([]string, 0, len(input.Access))
		accessList := make([]access.Access, 0, len(input.Access))
		for _, key := range input.Access {
			id, err := access.Parse(key)
			if err != nil {
				return fmt.Errorf("access %s doesn't exist. ", key)
			}
			accessIDs = append(accessIDs, strconv.Itoa(int(id)))
			accessList = append(accessList, id)
		}

		roleAccess, err := u.roleAccessStore.GetByRoleID(txCtx, role.ID)
		if err != nil {
			return err
		}

		//保存旧的accessIds
		oldAccessIDs := strings.Split(roleAccess.AccessIDs, ",")

		roleAccess.AccessIDs = strings.Join(accessIDs, ",")
		if err = u.roleAccessStore.Save(txCtx, roleAccess); err != nil {
			return err
		}
		//更新角色缓存
		roleAccessCacheConf := model.CreateRoleAccess(accessList...)

		if err = u.roleAccessCache.Set(ctx, u.roleAccessCache.Key(roleUUID), roleAccessCacheConf, time.Hour); err != nil {
			return err
		}

		//角色权限变更日志
		return u.roleAccessLogStore.HistoryEdit(txCtx, 0, role.ID, &entry.AccessListLog{
			AccessIds: oldAccessIDs,
			Module:    role.Module,
		}, &entry.AccessListLog{
			AccessIds: accessIDs,
			Module:    role.Module,
		}, userID)
	})
}

func (u *userInfoService) DeleteRole(ctx context.Context, userID int, roleUUID string) error {
	//内置角色不可以修改
	if access.IsBuildInRole(roleUUID) {
		return errors.New("build-in role can't be deleted. ")
	}

	//判断角色uuid存不存在
	role, err := u.roleStore.GetByUUID(ctx, roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("role_id is invalid. ")
		}
		return err
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: roleUUID,
		Name: role.Title,
	})

	t := time.Now()
	return u.roleStore.Transaction(ctx, func(txCtx context.Context) error {
		//从角色表中删除角色
		if _, err = u.roleStore.Delete(txCtx, role.ID); err != nil {
			return err
		}
		delMap := make(map[string]interface{})
		delMap["`role_id`"] = role.ID

		roleAccess, err := u.roleAccessStore.GetByRoleID(txCtx, role.ID)
		if err != nil {
			return err
		}

		//角色权限表中删除 角色记录
		if _, err = u.roleAccessStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		//从用户角色表中，获取有这角色的所有用户
		userRoleList, err := u.userRoleStore.GetListByRoleID(txCtx, role.ID)
		if err != nil {
			return err
		}
		//更新用户表中用户的role_ids
		for _, userRole := range userRoleList {
			user, err := u.userInfoStore.Get(txCtx, userRole.UserID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue
				}
				return err
			}

			roleIds := strings.Split(user.RoleIds, ",")
			idx := slices.Index(roleIds, roleUUID)
			//若用户没有此角色则跳过
			if idx == -1 {
				continue
			}
			user.RoleIds = strings.Join(slices.Delete(roleIds, idx, idx+1), ",")
			user.UpdateTime = t
			_, err = u.userInfoStore.Update(txCtx, user)
			if err != nil {
				return err
			}
		}

		//用户角色表中 删除对应角色的所有用户
		if _, err = u.userRoleStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		//角色权限变更日志
		err = u.roleAccessLogStore.HistoryDelete(txCtx, 0, role.ID, &entry.AccessListLog{
			AccessIds: strings.Split(roleAccess.AccessIDs, "m"),
			Module:    role.Module,
		}, userID)

		//更新角色缓存
		return u.roleAccessCache.Delete(ctx, u.roleAccessCache.Key(roleUUID))
	})
}

func (u *userInfoService) RoleBatchUpdate(ctx context.Context, userIds []int, roleUUID string) error {
	roleID := 0
	if access.IsBuildInRole(roleUUID) {
		buildInRole := access.GetBuildInRole(roleUUID)
		if !buildInRole.IsAddable {
			return fmt.Errorf("build-in role %s is unaddable. ", roleUUID)
		}
		roleID = buildInRole.ID
	} else {
		roleInfo, err := u.roleStore.GetByUUID(ctx, roleUUID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("role doesn't exist. ")
			}
			return err
		}
		roleID = roleInfo.ID
	}

	userIdsUpdateList := make([]int, 0, len(userIds))
	userInfoUpdateList := make([]*entry.UserInfo, 0, len(userIds))
	err := u.roleStore.Transaction(ctx, func(txCtx context.Context) error {

		t := time.Now()
		for _, userID := range userIds {
			user, err := u.userInfoStore.Get(txCtx, userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue
				}
				return err
			}
			//若用户已有此角色
			roleIds := strings.Split(user.RoleIds, ",")
			idx := slices.Index(roleIds, roleUUID)
			if idx != -1 {
				continue
			}

			/*用户角色表中 删除用户原有的角色
			TODO 现在是单角色，需要删除原有角色 之后多角色后，不需要删除已有角色，但要校验用户是否已有这个角色
			*/
			delMap := make(map[string]interface{})
			delMap["`user_id`"] = user.Id
			if _, err = u.userRoleStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}
			userRole := &entry.UserRole{
				UserID:     user.Id,
				RoleID:     roleID,
				Module:     "/",
				CreateTime: t,
			}
			err = u.userRoleStore.Insert(txCtx, userRole)
			if err != nil {
				return err
			}

			//更新用户表role_ids, TODO 现在是单角色，直接赋值
			//roleIds = append(roleIds, roleUUID)
			user.RoleIds = roleUUID
			user.UpdateTime = t
			_, err = u.userInfoStore.Update(txCtx, user)
			if err != nil {
				return err
			}

			userIdsUpdateList = append(userIdsUpdateList, user.Id)
			userInfoUpdateList = append(userInfoUpdateList, user)
		}
		return nil
	})
	if err != nil {
		return err
	}

	userIdStrs := make([]string, 0)
	userNames := make([]string, 0)
	//等事务执行完才刷新缓存，防止批量操作中已刷新缓存但事务失败的情况
	for i, userID := range userIdsUpdateList {
		//更新用户缓存
		_ = u.userInfoCache.Set(ctx, u.userInfoCache.Key(userID), userInfoUpdateList[i], time.Minute*30)

		userIdStrs = append(userIdStrs, strconv.Itoa(userID))
		userNames = append(userNames, userInfoUpdateList[i].UserName)

	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{Uuid: strings.Join(userIdStrs, ","), Name: strings.Join(userNames, ",")})
	return nil
}

func (u *userInfoService) RoleBatchRemove(ctx context.Context, userIds []int, roleUUID string) error {
	//如果不是内置角色，判断该自定义角色存不存在
	if !access.IsBuildInRole(roleUUID) {
		_, err := u.roleStore.GetByUUID(ctx, roleUUID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("role doesn't exist. ")
			}
			return err
		}
	}

	userIdsUpdateList := make([]int, 0, len(userIds))
	userInfoUpdateList := make([]*entry.UserInfo, 0, len(userIds))
	err := u.roleStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		for _, userID := range userIds {
			user, err := u.userInfoStore.Get(txCtx, userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue
				}
				return err
			}
			//若用户没有此角色则跳过
			roleIds := strings.Split(user.RoleIds, ",")
			idx := slices.Index(roleIds, roleUUID)
			if idx == -1 {
				continue
			}
			//用户角色表中 删除用户原有的角色
			delMap := make(map[string]interface{})
			delMap["`user_id`"] = user.Id
			if _, err = u.userRoleStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}
			//更新用户表role_ids
			user.RoleIds = strings.Join(slices.Delete(roleIds, idx, idx+1), ",")
			user.UpdateTime = t
			err = u.userInfoStore.Save(txCtx, user)
			if err != nil {
				return err
			}

			userIdsUpdateList = append(userIdsUpdateList, user.Id)
			userInfoUpdateList = append(userInfoUpdateList, user)
		}

		return nil
	})
	if err != nil {
		return err
	}

	userIdStrs := make([]string, 0)
	userNames := make([]string, 0)
	//等事务执行完才刷新缓存，防止批量操作中已刷新缓存但事务失败的情况
	for i, userID := range userIdsUpdateList {
		//更新用户缓存
		_ = u.userInfoCache.Set(ctx, u.userInfoCache.Key(userID), userInfoUpdateList[i], time.Minute*30)

		userIdStrs = append(userIdStrs, strconv.Itoa(userID))
		userNames = append(userNames, userInfoUpdateList[i].UserName)
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{Uuid: strings.Join(userIdStrs, ","), Name: strings.Join(userNames, ",")})

	return nil
}

func (u *userInfoService) GetUserRoleIds(ctx context.Context, userID int) ([]string, error) {
	userInfo, err := u.userInfoStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return strings.Split(userInfo.RoleIds, ","), nil
}

func (u *userInfoService) GetRoleAccessIds(ctx context.Context, roleUUID string) ([]string, error) {
	if access.IsBuildInRole(roleUUID) {
		role := access.GetBuildInRole(roleUUID)
		accessIDs := make([]string, len(role.AccessID))
		for i, ac := range role.AccessID {
			accessIDs[i] = strconv.Itoa(int(ac))
		}

		return accessIDs, nil
	}

	role, err := u.roleStore.GetByUUID(ctx, roleUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, err
	}
	roleAccess, err := u.roleAccessStore.GetByRoleID(ctx, role.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, err
	}
	return strings.Split(roleAccess.AccessIDs, ","), nil
}
