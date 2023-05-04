package audit_controller

import (
	"bytes"
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	apintoModule "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

func (a *auditLogController) Handler(ginCtx *gin.Context) {
	start := time.Now()
	bodyReader := ginCtx.Request.Body
	bodyBytes, err := io.ReadAll(bodyReader)
	bodyReader.Close()
	if err != nil {
		log.Warn("read body :", err)
		ginCtx.JSON(http.StatusOK, controller.NewNoAccessError("Invalid request body"))
		return
	}
	ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	//将请求体加入到上下文中，兼容需要修改的情况， 比如批量上线接口
	ginCtx.Set(controller.LogBody, string(bodyBytes))

	ginCtx.Next()

	kind := ginCtx.GetString(apintoModule.ApintoModuleName)
	operate := ginCtx.GetInt(controller.Operate)
	if operate == 0 {
		operate = switchMethod(ginCtx.Request.Method)
	}
	end := time.Now()
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := users.GetUserId(ginCtx)

	//获取ip
	ip := ""
	if forwarded := ginCtx.GetHeader("X-Forwarded-For"); len(forwarded) > 0 {
		if i := strings.Index(forwarded, ","); i > 0 {
			ip = forwarded[:i]
		} else {
			ip = forwarded
		}
	} else if realIP := ginCtx.GetHeader("X-Real-Ip"); len(realIP) > 0 {
		ip = realIP
	} else {
		ip = ginCtx.ClientIP()
	}

	userAgent := ginCtx.GetHeader("user-agent")

	errInfo := ginCtx.GetString(controller.ErrorMessage)

	url := fmt.Sprintf("%s %s", ginCtx.Request.Method, ginCtx.Request.RequestURI)

	a.auditLogService.Log(namespaceId, userId, operate, kind, url, ginCtx.GetString("auditObject"), ip, userAgent, ginCtx.GetString(controller.LogBody), errInfo, start, end)

}
func switchMethod(method string) int {
	switch method {
	case http.MethodGet:
		return int(audit_model.LogOperateTypeNone)
	case http.MethodPost:
		return int(audit_model.LogOperateTypeCreate)
	case http.MethodPut, http.MethodPatch:
		return int(audit_model.LogOperateTypeEdit)
	case http.MethodDelete:
		return int(audit_model.LogOperateTypeDelete)
	}
	return int(audit_model.LogOperateTypeNone)
}
