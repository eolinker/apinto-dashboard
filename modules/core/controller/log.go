package controller

import (
	"bytes"
	"context"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
)

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
