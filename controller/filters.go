package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	UserId         = "userId"
	Authorization  = "Authorization"
	Session        = "Session"
	RAuthorization = "RAuthorization"
)

var loginError = errors.New("登录已失效，请重新登录")

type IFilterRouter interface {
	MustNamespace(*gin.Context)
	VerifyToken(*gin.Context)
}

func VerifyToken(ginCtx *gin.Context) {

	session, _ := ginCtx.Cookie(Session)
	if session == "" {
		session = uuid.New()
		ginCtx.SetCookie(Session, session, 0, "", "", false, true)

	}

	ginCtx.Set(UserId, 1)
}

func GetUserId(ginCtx *gin.Context) int {
	return ginCtx.GetInt(UserId)
}

func GenAccessHandler(acs ...access.Access) gin.HandlerFunc {

	return func(ginCtx *gin.Context) {
		// todo 原实现不适合开源，这里埋点用于以后扩展
	}
}

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic %v", common.PanicTrace(err))
			fmt.Println(common.PanicTrace(err))
			c.JSON(http.StatusInternalServerError, NewErrorResult("服务器内部错误"))
		}
	}()

	c.Next()
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(ginCtx *gin.Context) {
	if !strings.HasPrefix(ginCtx.Request.URL.Path, "/api/") {
		ginCtx.Next()
		return
	}

	ginCtx.Request = ginCtx.Request.WithContext(context.Background())

	data, err := ginCtx.GetRawData()
	if err == nil {
		ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(data)) // 关键点
	}

	blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
	ginCtx.Writer = blw

	// 开始时间
	startTime := time.Now()

	// 处理请求
	ginCtx.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := ginCtx.Request.Method

	// 请求路由
	reqUri := ginCtx.Request.RequestURI

	// 状态码
	statusCode := ginCtx.Writer.Status()

	// 请求IP
	clientIP := ginCtx.ClientIP()
	resBody := blw.Body.String()
	//日志格式
	log.DebugF("| %3d | %13v | %15s | %s | %s | reqBody=%s | resBody=%s ",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri,
		string(data),
		resBody,
	)
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
