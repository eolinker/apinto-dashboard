package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ICore interface {
	http.Handler
	ResetVersion(version string)
	ReloadModule() error
	CheckNewModule(pluginID, name, apiGroup string, config interface{}) error
}

type EngineCreate interface {
	CreateEngine() (engine *gin.Engine)
}
