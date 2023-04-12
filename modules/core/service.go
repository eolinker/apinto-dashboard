package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ICore interface {
	http.Handler
	ResetVersion(version string)
	ReloadModule() error
	CheckNewModule(uuid, name, driver string, define, config interface{}) error
	HasModule(module string, path string) bool
}

type EngineCreate interface {
	CreateEngine() (engine *gin.Engine)
}
