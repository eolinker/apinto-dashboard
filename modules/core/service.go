package core

import (
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ICore interface {
	http.Handler
	ResetVersion(version string)
	ReloadModule() error
	CheckNewModule(pluginID, name, apiGroup string, config interface{}) error
}

type IProviders interface {
	apinto_module.IProviders
	Set(providers apinto_module.IProviders)
}
type EngineCreate interface {
	CreateEngine() (engine *gin.Engine)
}
