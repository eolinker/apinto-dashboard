package user

import (
	"context"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	user_model2 "github.com/eolinker/apinto-dashboard/modules/user/user-model"
)

type IUserInfoService interface {
	GetUserInfo(ctx context.Context, userId int) (*user_model2.UserInfo, error)
	UpdateMyProfile(ctx context.Context, userId int, req *user_dto.UpdateMyProfileReq) error
	UpdateMyPassword(ctx context.Context, userId int, req *user_dto.UpdateMyPasswordReq) error
	UpdateLoginTime(ctx context.Context, userId int)
	GetUserInfoList(ctx context.Context, roleId, keyWord string) ([]*user_model2.UserInfo, error)
	GetUserInfoAll(ctx context.Context) ([]*user_model2.UserInfo, error)
	GetUserInfoMaps(ctx context.Context, userId ...int) (map[int]*user_model2.UserInfo, error)
	CreateUser(ctx context.Context, operator int, userInfo *user_dto.SaveUserReq) error
	CheckUser(ctx context.Context, userId int) error
	PatchUser(ctx context.Context, operator, userId int, req *user_dto.PatchUserReq) error
	UpdateUser(ctx context.Context, operator, userId int, userInfo *user_dto.SaveUserReq) error
	DelUser(ctx context.Context, operator int, userIds []int) error
	ResetUserPwd(ctx context.Context, operator, userId int, pwd string) error

	GetAccessInfo(ctx context.Context, userId int) (map[access.Access]struct{}, error)
	GetRoleList(ctx context.Context, userID int) ([]*user_model2.RoleListItem, int, error)
	GetRoleInfo(ctx context.Context, roleID string) (*user_model2.RoleInfo, error)
	GetRoleOptions(ctx context.Context) ([]*user_model2.RoleOptionItem, error)
	CreateRole(ctx context.Context, userID int, input *user_dto.ProxyRoleInfo) error
	UpdateRole(ctx context.Context, userID int, roleUUID string, input *user_dto.ProxyRoleInfo) error
	DeleteRole(ctx context.Context, userID int, roleUUID string) error
	RoleBatchUpdate(ctx context.Context, userIds []int, roleUUID string) error
	RoleBatchRemove(ctx context.Context, userIds []int, roleUUID string) error

	GetUserRoleIds(ctx context.Context, userID int) ([]string, error)
	GetRoleAccessIds(ctx context.Context, roleUUID string) ([]string, error)
	CreateAdmin() error
	CleanAdminCache()
}
