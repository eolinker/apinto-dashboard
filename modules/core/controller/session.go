package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var jwtSecret = []byte("apintp-dashboard")

const (
	loginError = "登录已失效，请重新登录"
)

func (u *UserController) LoginCheckApi(ginCtx *gin.Context) {

	session, _ := ginCtx.Cookie(controller.Session)
	if session == "" {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		ginCtx.Abort()
		return
	}

	tokens, err := u.sessionCache.Get(ginCtx, session)
	if err == redis.Nil || tokens == nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		ginCtx.Abort()
		return
	}

	token := tokens.Jwt
	rToken := tokens.RJwt

	uc, err := common.JWTDecode(token, jwtSecret)
	if err != nil {
		controller.ErrorJsonWithCode(ginCtx, http.StatusOK, controller.CodeLoginInvalid, loginError)
		ginCtx.Abort()
		return
	}

	ginCtx.Writer.Header().Set(controller.Authorization, token)
	ginCtx.Writer.Header().Set(controller.RAuthorization, rToken)

	ginCtx.Set(controller.UserName, uc.Uname)
	ginCtx.Set(controller.Session, session)
}
func (u *UserController) LoginCheckModule(ginCtx *gin.Context) {

	session, _ := ginCtx.Cookie(controller.Session)
	if session == "" {

		doRedirect(ginCtx)
		return
	}

	tokens, err := u.sessionCache.Get(ginCtx, session)
	if err == redis.Nil || tokens == nil {
		doRedirect(ginCtx)
		return
	}

	token := tokens.Jwt
	rToken := tokens.RJwt

	uc, err := common.JWTDecode(token, jwtSecret)
	if err != nil {
		doRedirect(ginCtx)
		return
	}

	ginCtx.Writer.Header().Set(controller.Authorization, token)
	ginCtx.Writer.Header().Set(controller.RAuthorization, rToken)

	ginCtx.Set(controller.UserName, uc.Uname)
	ginCtx.Set(controller.Session, session)
}

func doRedirect(ginCtx *gin.Context) {

	url := fmt.Sprint("/login", "?callback=", ginCtx.Request.RequestURI)
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Redirect(http.StatusFound, url)
	ginCtx.Abort()
}
func (u *UserController) SetUser(ginCtx *gin.Context) {
	userId := users.GetUserId(ginCtx)
	if userId != 0 {
		return
	}
	cookie, err := ginCtx.Cookie(controller.Session)
	if err != nil {
		return
	}
	if ginCtx.GetString(controller.Session) != "" {
		return
	}
	session, _ := u.sessionCache.Get(ginCtx, cookie)
	if session == nil {
		return
	}
	uc := controller.UserClaim{}
	parse, err := jwt.ParseWithClaims(session.Jwt, &uc, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return
	}
	if parse.Valid {
		ginCtx.Set(controller.UserName, uc.Uname)
	}
}
