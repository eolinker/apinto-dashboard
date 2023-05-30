package controller

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func SetGinContextAuditObject(ctx context.Context, v interface{}) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		objectData, _ := json.Marshal(v)
		ginCtx.Set(AuditObject, string(objectData))

	}

}

func GetGinContextAuditObject(ginCtx *gin.Context) (value any, exists bool) {

	return ginCtx.Get(AuditObject)

}

func GetNamespaceId(ginCtx *gin.Context) int {
	id := ginCtx.GetInt(NamespaceId)
	if id == 0 {
		return 1
	}
	return id
}
