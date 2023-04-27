package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/user"
	user_dto "github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userInfo user.IUserInfoService
}

func (u *UserController) myProfile(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)

	userInfo, err := u.userInfo.GetUserInfo(ginCtx, userId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("getMyProfile fail. err:%s", err.Error()))

		return
	}

	lastLogin := ""
	if userInfo.LastLoginTime != nil {
		lastLogin = common.TimeToStr(*userInfo.LastLoginTime)
	}
	resUserInfo := user_dto.UserInfo{
		Id:     userInfo.Id,
		Sex:    userInfo.Sex,
		Avatar: userInfo.Avatar,

		Email: userInfo.Email,
		Phone: userInfo.Phone,

		UserName:     userInfo.UserName,
		NickName:     userInfo.NickName,
		NoticeUserId: userInfo.NoticeUserId,
		LastLogin:    lastLogin,
	}

	m := make(map[string]interface{})
	m["profile"] = resUserInfo
	m["describe"] = ""

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
func (u *UserController) myProfileUpdate(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)

	req := &user_dto.UpdateMyProfileReq{}
	err := ginCtx.BindJSON(req)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("updateMyProfile fail. err:%s", err.Error()))
		return
	}

	if err = u.userInfo.UpdateMyProfile(ginCtx, userId, req); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("updateMyProfile fail. err:%s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
func (u *UserController) myPassword(ginCtx *gin.Context) {

}
func (u *UserController) userEnum(ginCtx *gin.Context) {

}
func (u *UserController) ssoLogin(ginCtx *gin.Context) {

}
func (u *UserController) ssoLogout(ginCtx *gin.Context) {

}
func (u *UserController) ssoLoginCheck(ginCtx *gin.Context) {

}
func newUserController() *UserController {
	u := &UserController{}
	bean.Autowired(&u.userInfo)
	return u
}

func userRouters() apinto_module.RoutersInfo {
	userController := newUserController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/my/profile",
			Handler:     "core.my.profile",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.myProfile},
			Alternative: true,
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/my/profile",
			Handler:     "core.my.profile.reset",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.myProfileUpdate},
			Alternative: true,
		}, {
			Method:      http.MethodPost,
			Path:        "/api/my/password",
			Handler:     "core.my.password",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.myPassword},
			Alternative: true,
		}, {
			Method:      http.MethodGet,
			Path:        "/api/user/enum",
			Handler:     "core.user.enum",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.userEnum},
			Alternative: true,
		}, {
			Method:      http.MethodPost,
			Path:        "/sso/login",
			Handler:     "core.sso.login",
			Labels:      apinto_module.RouterLabelAnonymous,
			HandlerFunc: []apinto_module.HandlerFunc{userController.ssoLogin},
			Alternative: true,
		},
		{
			Method:      http.MethodPost,
			Path:        "/sso/login/check",
			Handler:     "core.sso.login.check",
			Labels:      apinto_module.RouterLabelAnonymous,
			HandlerFunc: []apinto_module.HandlerFunc{userController.ssoLoginCheck},
			Alternative: true,
		}, {
			Method:      http.MethodPost,
			Path:        "/sso/logout",
			Handler:     "core.sso.logout",
			Labels:      apinto_module.RouterLabelAnonymous,
			HandlerFunc: []apinto_module.HandlerFunc{userController.ssoLogout},
			Alternative: true,
		},
	}
}
