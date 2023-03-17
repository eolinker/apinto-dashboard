package user_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type userController struct {
	userInfo user.IUserInfoService
}

func RegisterUserRouter(router gin.IRoutes) {
	u := &userController{}
	bean.Autowired(&u.userInfo)

	router.GET("/access", u.getAllAccess)
	router.GET("/my/modules", u.getUserAccess)
	router.GET("/my/profile", u.getMyProfile)

}

func (u *userController) getMyProfile(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)

	userInfo, err := u.userInfo.GetUserInfo(ginCtx, userId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("getMyProfile fail. err:%s", err.Error())))
		return
	}

	lastLogin := ""
	if userInfo.LastLoginTime != nil {
		lastLogin = common.TimeToStr(*userInfo.LastLoginTime)
	}
	resUserInfo := user_dto.UserInfo{
		Id:           userInfo.Id,
		Sex:          userInfo.Sex,
		Avatar:       userInfo.Avatar,
		Desc:         userInfo.Remark,
		Email:        userInfo.Email,
		Phone:        userInfo.Phone,
		Status:       userInfo.Status,
		UserName:     userInfo.UserName,
		NickName:     userInfo.NickName,
		NoticeUserId: userInfo.NoticeUserId,
		LastLogin:    lastLogin,
		CreateTime:   common.TimeToStr(userInfo.CreateTime),
		UpdateTime:   common.TimeToStr(userInfo.UpdateTime),
		RoleIds:      strings.Split(userInfo.RoleIds, ","),
	}

	m := make(map[string]interface{})
	m["profile"] = resUserInfo
	m["describe"] = userInfo.Remark

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

func (u *userController) getAllAccess(ginCtx *gin.Context) {

	globalAccess, depth := access.GetGlobalAccessConfig()

	data := make(map[string]interface{})
	data["modules"] = getModules(globalAccess)
	data["depth"] = depth

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))

}

func getModules(modules []*access.GlobalAccess) []*user_dto.SystemModuleItem {
	items := make([]*user_dto.SystemModuleItem, len(modules))
	for i, module := range modules {
		item := &user_dto.SystemModuleItem{
			ID:     module.ID,
			Title:  module.Title,
			Module: module.Module,
			Access: module.Access,
		}
		if len(module.Children) > 0 {
			item.Children = getModules(module.Children)
		}
		items[i] = item
	}
	return items
}

func (u *userController) getUserAccess(ginCtx *gin.Context) {
	userID := controller.GetUserId(ginCtx)

	accessSet, err := u.userInfo.GetAccessInfo(ginCtx, userID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("GetUserAccessList fail. err:%s", err.Error())))
		return
	}
	allModules := access.GetAllModulesConfig()

	modules := make([]*user_dto.UserModuleItem, 0, len(accessSet))
	for _, module := range allModules {
		for _, needId := range module.ModuleNeed {
			if _, has := accessSet[needId]; has {
				accessList := make([]string, 0, len(module.Access))
				for _, key := range module.Access {
					accessId, err := access.Parse(key)
					if err != nil {
						continue
					}
					if _, exist := accessSet[accessId]; exist {
						accessList = append(accessList, key)
					}
				}
				modules = append(modules, &user_dto.UserModuleItem{
					Id:     module.ID,
					Router: module.Router,
					Title:  module.Title,
					Access: accessList,
					Parent: module.Parent,
				})
				break
			}
		}
	}

	data := make(map[string]interface{})
	data["modules"] = modules

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))

}
