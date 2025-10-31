package apinto_module

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
)

type Module = pm3.Module
type ModuleNeedKill interface {
	Kill()
}

type HandlerFunc = gin.HandlerFunc

type RouterInfo = pm3.Api

type RoutersInfo = []pm3.Api

type MiddlewareHandler = pm3.Middleware
