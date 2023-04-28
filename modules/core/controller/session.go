package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
		return
	}

	ginCtx.Writer.Header().Set(controller.Authorization, token)
	ginCtx.Writer.Header().Set(controller.RAuthorization, rToken)

	userId := uc.Id

	ginCtx.Set(controller.UserId, userId)
	ginCtx.Set(controller.Session, session)
}
func (u *UserController) SetUser(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)
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
		ginCtx.Set(controller.UserId, uc.Id)
	}
}
