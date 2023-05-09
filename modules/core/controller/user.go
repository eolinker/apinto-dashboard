package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller/users"
	random_controller "github.com/eolinker/apinto-dashboard/modules/base/random-controller"
	"net/http"
	"time"

	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"

	"github.com/eolinker/apinto-dashboard/cache"

	"github.com/go-basic/uuid"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/user"
	user_dto "github.com/eolinker/apinto-dashboard/modules/user/user-dto"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userInfo      user.IUserInfoService
	sessionCache  user.ISessionCache
	commonCache   cache.ICommonCache
	moduleService module_plugin.IModulePlugin
}

func (u *UserController) myProfile(ginCtx *gin.Context) {
	userId := users.GetUserId(ginCtx)

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
	userId := users.GetUserId(ginCtx)

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
func (u *UserController) setPassword(ginCtx *gin.Context) {
	userId := users.GetUserId(ginCtx)

	req := &user_dto.UpdateMyPasswordReq{}
	err := ginCtx.BindJSON(req)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("updateMyPassword fail. err:%s", err.Error()))
		return
	}

	if err = u.userInfo.UpdateMyPassword(ginCtx, userId, req); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("updateMyPassword fail. err:%s", err.Error()))
		return

	}
	info, err := u.userInfo.GetUserInfo(ginCtx, userId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("get user info fail. err:%s", err.Error()))
		return
	}
	userJWT, err := common.JWTEncode(&controller.UserClaim{
		Id:        info.Id,
		Uname:     info.UserName,
		LoginTime: info.LastLoginTime.Format("2006-01-02 15:04:05"),
	}, jwtSecret)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, "登录失效")
		return
	}
	cookieValue := uuid.New()
	session := &user_model.Session{
		Jwt:  userJWT,
		RJwt: common.Md5(userJWT),
	}

	if err = u.sessionCache.Set(ginCtx, cookieValue, session, time.Hour*24*7); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//每次登陆都把之前的token清除掉
	{
		userIdCookieKey := fmt.Sprintf("userId:%d", info.Id)
		oldCookie, _ := u.commonCache.Get(ginCtx, userIdCookieKey)
		_ = u.sessionCache.Delete(ginCtx, string(oldCookie))
		_ = u.commonCache.Set(ginCtx, userIdCookieKey, []byte(cookieValue), time.Hour*24*7)
	}

	ginCtx.SetCookie(controller.Session, cookieValue, 0, "", "", false, true)

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}
func (u *UserController) userEnum(ginCtx *gin.Context) {
	userInfoList, err := u.userInfo.GetAllUsers(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("获取用户列表失败. err:%s", err.Error()))
		return
	}

	resList := make([]user_dto.UserInfo, 0, len(userInfoList))

	for _, userInfo := range userInfoList {
		lastLogin := ""
		if userInfo.LastLoginTime != nil {
			lastLogin = common.TimeToStr(*userInfo.LastLoginTime)
		}
		resUserInfo := user_dto.UserInfo{
			Id:           userInfo.Id,
			Sex:          userInfo.Sex,
			Avatar:       userInfo.Avatar,
			Email:        userInfo.Email,
			Phone:        userInfo.Phone,
			UserName:     userInfo.UserName,
			NickName:     userInfo.NickName,
			NoticeUserId: userInfo.NoticeUserId,
			LastLogin:    lastLogin,
		}
		resList = append(resList, resUserInfo)
	}

	m := make(map[string]interface{})
	m["users"] = resList

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
func (u *UserController) access(ginCtx *gin.Context) {
	modules, err := u.moduleService.GetEnabledPlugins(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("get module info fail. err:%s", err.Error()))
		return
	}
	access := make([]*user_dto.UserAccess, 0, len(modules))
	for _, m := range modules {
		access = append(access, &user_dto.UserAccess{
			Name:   m.Name,
			Access: "edit",
		})
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"access": access,
	}))
}

func (u *UserController) ssoLogin(ginCtx *gin.Context) {
	var loginInfo user_dto.UserLogin
	err := ginCtx.BindJSON(&loginInfo)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("bind login param fail. err:%s", err.Error()))
		return
	}

	id, success := u.userInfo.CheckPassword(ginCtx, loginInfo.Username, loginInfo.Password)
	if !success {
		ginCtx.JSON(http.StatusOK, controller.NewLoginInvalidError(controller.CodeLoginPwdErr, "登录失败，用户名或密码错误"))
		return
	}

	now := time.Now()
	// 成功登录，更新登录时间
	err = u.userInfo.UpdateLastLoginTime(ginCtx, id, &now)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, "登录失败")
		return
	}

	userJWT, err := common.JWTEncode(&controller.UserClaim{
		Id:        id,
		Uname:     loginInfo.Username,
		LoginTime: now.Format("2006-01-02 15:04:05"),
	}, jwtSecret)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, "登录失败")
		return
	}
	cookieValue := uuid.New()
	session := &user_model.Session{
		Jwt:  userJWT,
		RJwt: common.Md5(userJWT),
	}

	if err = u.sessionCache.Set(ginCtx, cookieValue, session, time.Hour*24*7); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//每次登陆都把之前的token清除掉
	{
		userIdCookieKey := fmt.Sprintf("userId:%d", id)
		oldCookie, _ := u.commonCache.Get(ginCtx, userIdCookieKey)
		_ = u.sessionCache.Delete(ginCtx, string(oldCookie))
		_ = u.commonCache.Set(ginCtx, userIdCookieKey, []byte(cookieValue), time.Hour*24*7)
	}

	ginCtx.SetCookie(controller.Session, cookieValue, 0, "", "", false, true)
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
func (u *UserController) ssoLogout(ginCtx *gin.Context) {
	cookie, err := ginCtx.Cookie(controller.Session)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	u.sessionCache.Delete(ginCtx, cookie)
	defer func() {
		ginCtx.SetCookie(controller.Session, cookie, -1, "", "", false, true)
	}()

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
func (u *UserController) ssoLoginCheck(ginCtx *gin.Context) {
	cookie, err := ginCtx.Cookie(controller.Session)
	if err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		return
	}

	session, _ := u.sessionCache.Get(ginCtx, cookie)
	if session == nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		return
	}
	uc, err := common.JWTDecode(session.Jwt, jwtSecret)
	if err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		return
	}

	info, err := u.userInfo.GetUserInfo(ginCtx, uc.Id)
	if err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		return
	}
	if info.LastLoginTime.Format("2006-01-02 15:04:05") != uc.LoginTime || info.UserName != uc.Uname {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(true))
}
func newUserController() *UserController {
	u := &UserController{}
	bean.Autowired(&u.userInfo)
	bean.Autowired(&u.sessionCache)
	bean.Autowired(&u.commonCache)
	bean.Autowired(&u.moduleService)
	return u
}
func randomRouters() apinto_module.RoutersInfo {
	r := random_controller.NewRandomController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/random/:template/id",
			Handler:     "core.random.id",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{r.GET},
			Alternative: false,
		}}
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
			HandlerFunc: []apinto_module.HandlerFunc{userController.setPassword},
			Alternative: true,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/my/access",
			Handler:     "core.my.access",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.access},
			Alternative: true,
		}, {
			Method:      http.MethodGet,
			Path:        "/api/user/enum",
			Handler:     "core.user.enum",
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{userController.userEnum},
			Alternative: true,
		},
		{
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
