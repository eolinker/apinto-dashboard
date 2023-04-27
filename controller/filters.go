package controller

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

const (
	UserId         = "userId"
	Authorization  = "Authorization"
	Session        = "Session"
	RAuthorization = "RAuthorization"
	ErrorCode      = "ErrorCode"
	ErrorMessage   = "ErrorBody"
	Operate        = "Operate"
	AuditKind      = "AuditKind"
	LogBody        = "LogBody"
	AuditObject    = "auditObject"
)

func GetUserId(ginCtx *gin.Context) int {
	return ginCtx.GetInt(UserId)
}

type UserClaim struct {
	Id        int    `json:"id"`
	Uname     string `json:"username"`
	LoginTime string `json:"login_time"`
	jwt.StandardClaims
}

func GenAccessHandler() gin.HandlerFunc {

	return func(ginCtx *gin.Context) {
		// todo 原实现不适合开源，这里埋点用于以后扩展
	}
}

func ErrorJsonWithCode(ginCtx *gin.Context, statusCode int, errorCode int, errorMsg string) {
	ginCtx.Set(ErrorMessage, errorMsg)
	ginCtx.Set(ErrorCode, errorCode)
	ginCtx.JSON(statusCode, &Result{
		Code: errorCode,
		Msg:  errorMsg,
	})

}
func ErrorJson(ginCtx *gin.Context, statusCode int, errorMsg string) {
	ErrorJsonWithCode(ginCtx, statusCode, ordinaryCode, errorMsg)
}
