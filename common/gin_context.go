package common

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func SetGinContextAuditObject(ctx context.Context, v interface{}) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		objectData, _ := json.Marshal(v)
		ginCtx.Set("auditObject", string(objectData))

	}

}

func GetGinContextAuditObject(ginCtx *gin.Context) (value any, exists bool) {

	return ginCtx.Get("auditObject")

}
