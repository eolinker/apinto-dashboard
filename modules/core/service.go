package core

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/core/model"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ICore interface {
	http.Handler
	ResetVersion(version string)
	ReloadModule() error
	CheckNewModule(uuid, name, driver string, define, config interface{}) error
	HasModule(module string, path string) bool
	SetCoreModule(module apinto_module.CoreModule)
}

type EngineCreate interface {
	CreateEngine() (engine *gin.Engine)
}
type INavigationService interface {
	List(ctx context.Context) ([]*model.Navigation, map[string]string, error)
}
