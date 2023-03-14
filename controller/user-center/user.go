package user_center

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/dto/user-dto"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/apinto-dashboard/user_center/client"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
	"strconv"
	"time"
)

var (
	userCenterClient client.IUserCenterClient
	userInfoService  service.IUserInfoService
	sessionCache     cache.ISessionCache
	commonCache      cache.ICommonCache
)

const exp = 3600 * 24 * 7

func RegisterUserCenterProxyRouter(router gin.IRouter) {

	bean.Autowired(&userCenterClient)
	bean.Autowired(&userInfoService)
	bean.Autowired(&sessionCache)
	bean.Autowired(&commonCache)

	router.POST("/sso/login", func(ginCtx *gin.Context) {
		input := new(user_dto.UserLoginInput)
		err := ginCtx.BindJSON(input)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		//key := common.Md5(input.Username)
		//
		//decodePwd, err := common.Base64Decode(input.Password)
		//if err != nil {
		//	ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		//	return
		//}
		//
		//pwdDecrypter := common.CBCDecrypter(decodePwd, []byte(key), iv)
		//
		//md5Pwd := common.Md5(string(pwdDecrypter))

		req := &client.UserLoginReq{
			Username:      input.Username,
			Password:      input.Password,
			Client:        input.Client,
			Type:          input.Type,
			CheckCodeType: input.CheckCodeType,
			CheckCode:     input.CheckCode,
			AppType:       input.AppType,
			Ref:           input.Ref,
			RefUrl:        input.RefUrl,
		}

		response, code, err := userCenterClient.Login(req)
		if err != nil {
			switch code {
			case 132000001:
				ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginUserNoExistent, err.Error()))
			case 132000002:
				ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginPwdErr, err.Error()))
			case 132000003:
				ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginCodeErr, err.Error()))
			default:
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			}
			return
		}

		token, err := common.VerifyToken(response.Jwt)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		ctx := context.Background()

		claims := token.Claims.(jwt.MapClaims)
		userIdStr := claims["userId"]
		userId, _ := strconv.Atoi(userIdStr.(string))

		//校验用户是否可以登录（禁用 删除都不可登录）

		if err = userInfoService.CheckUser(ctx, userId); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		//更新用户最新登录时间
		go userInfoService.UpdateLoginTime(ctx, userId)

		//session存到本地缓存
		cookieValue := uuid.New()
		session := &model.Session{
			Jwt:  response.Jwt,
			RJwt: response.RJWT,
		}

		if err = sessionCache.Set(ctx, sessionCache.Key(cookieValue), session, time.Hour*24*7); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		//每次登陆都把之前的token清除掉
		{
			userIdCookieKey := fmt.Sprintf("userId:%d", userId)
			oldCookie, _ := commonCache.Get(ctx, userIdCookieKey)
			_ = sessionCache.Delete(ctx, sessionCache.Key(string(oldCookie)))
			_ = commonCache.Set(ctx, userIdCookieKey, []byte(cookieValue), time.Hour*24*7)
		}

		ginCtx.SetCookie(controller.Session, cookieValue, 0, "", "", false, true)

		resData := user_dto.UserLoginData{
			Jwt:  response.Jwt,
			RJWT: response.RJWT,
			Type: response.Type,
		}

		m := make(map[string]interface{})
		m["data"] = resData

		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
	})
	router.POST("/sso/logout", func(ginCtx *gin.Context) {

		cookie, err := ginCtx.Cookie(controller.Session)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
			return
		}

		ctx := context.Background()
		sessionCache.Delete(ctx, sessionCache.Key(cookie))
		defer func() {
			ginCtx.SetCookie(controller.Session, cookie, -1, "", "", false, true)
		}()

		session, err := sessionCache.Get(ctx, sessionCache.Key(cookie))
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
			return
		}

		if err = userCenterClient.Logout(session.Jwt); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
			return
		}

		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
	})

	router.POST("/sso/login/check", func(ginCtx *gin.Context) {

		cookie, err := ginCtx.Cookie(controller.Session)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}

		session, _ := sessionCache.Get(context.Background(), sessionCache.Key(cookie))
		if session == nil {
			ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginCodeErr, "登录失效"))
			return
		}

		req := &client.LoginCheck{
			Authorization: session.Jwt,
		}
		response, err := userCenterClient.LoginCheck(req)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginCodeErr, err.Error()))
			return
		}

		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(response))
	})
}
