package controller

import (
	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Authorization  = "Authorization"
	Session        = "Session"
	RAuthorization = "RAuthorization"
	ErrorCode      = "ErrorCode"
	ErrorMessage   = "ErrorBody"
	//Operate        = "Operate"
	AuditKind   = "AuditKind"
	LogBody     = "LogBody"
	AuditObject = "auditObject"
	UserName    = "userName"
	NamespaceId = "namespaceId"
)

type UserClaim struct {
	Id        int    `json:"id"`
	Uname     string `json:"username"`
	LoginTime string `json:"login_time"`
	jwt.RegisteredClaims
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
	ErrorJsonWithCode(ginCtx, statusCode, OrdinaryCode, errorMsg)
}
