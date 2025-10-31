package middleware

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
)

type handler struct {
	pm3.MiddlewareChecker
	handlerFunc gin.HandlerFunc
}

func Create(handlerFunc gin.HandlerFunc, checker pm3.MiddlewareChecker) pm3.Middleware {
	return &handler{MiddlewareChecker: checker, handlerFunc: handlerFunc}
}

func CreateF(handlerFunc gin.HandlerFunc, checkerFunc func(api pm3.ApiInfo) bool) pm3.Middleware {
	return Create(handlerFunc, CheckHandleFunc(checkerFunc))
}

func (h *handler) Handle(ginCtx *gin.Context) {
	h.handlerFunc(ginCtx)
}

type CheckHandleFunc func(api pm3.ApiInfo) bool

func (f CheckHandleFunc) Check(api pm3.ApiInfo) bool {
	return f(api)
}
