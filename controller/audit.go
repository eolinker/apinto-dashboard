package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	auditservice "github.com/eolinker/apinto-dashboard/modules/audit"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	controller        = newController()
	_          IAudit = (*audit)(nil)
)

type IAudit interface {
	Handler(operate int, kind string) gin.HandlerFunc
}

type audit struct {
	service auditservice.IAuditLogService
}

func (a *audit) Handler(operate int, kind string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		start := time.Now()
		bodyReader := ginCtx.Request.Body
		bodyBytes, err := io.ReadAll(bodyReader)
		bodyReader.Close()
		if err != nil {
			log.Warn("read body :", err)
			ginCtx.JSON(http.StatusOK, dto.NewNoAccessError("Invalid request body"))
			return
		}
		ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		//将请求体加入到上下文中，兼容需要修改的情况， 比如批量上线接口
		ginCtx.Set("logBody", string(bodyBytes))

		ginCtx.Next()

		//特殊情况. 通用分组controller暂只支持api分组写入审计日志
		if kind == enum.LogKindCommonGroup {
			groupType := ginCtx.Param("group_type")
			if groupType != "api" {
				return
			}
			kind = enum.LogKindAPIGroup
		}

		end := time.Now()
		namespaceId := namespace_controller.GetNamespaceId(ginCtx)
		userId := ginCtx.GetInt(UserId)

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

		errInfo := ""
		//logger中间件将*ResponseWriterWrapper赋值给了ginCtx.Writer，这里是可断言的
		blw, ok := ginCtx.Writer.(*ResponseWriterWrapper)
		if ok {
			errBody := blw.Body.String()
			result := new(dto.Result)
			_ = json.Unmarshal([]byte(errBody), result)
			if result.Code != 0 {
				errInfo = result.Msg
			}
		}

		url := fmt.Sprintf("%s %s", ginCtx.Request.Method, ginCtx.Request.RequestURI)

		a.service.Log(namespaceId, userId, operate, kind, url, ginCtx.GetString("auditObject"), ip, userAgent, ginCtx.GetString("logBody"), errInfo, start, end)
	}
}

func newController() *audit {

	c := &audit{}
	bean.Autowired(&c.service)
	return c
}

func LogHandler(operate int, kind string) gin.HandlerFunc {
	return controller.Handler(operate, kind)
}
