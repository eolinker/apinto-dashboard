package controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/gin-gonic/gin"
)

type LogOperateType audit_model.LogOperateType

func (l LogOperateType) Handler(ginCtx *gin.Context) {
	ginCtx.Set(Operate, l)
}
