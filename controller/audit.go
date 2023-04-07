package controller

import "github.com/gin-gonic/gin"

func AuditLogHandler(operate int, kind string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(AuditKind, kind)
		ctx.Set(Operate, operate)

		handlerFunc(ctx)
	}
}
