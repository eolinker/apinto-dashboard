package core

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/core/model"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ICore interface {
	http.Handler
	ResetVersion(version string)
	ReloadModule() error
}

type EngineCreate interface {
	CreateEngine() (engine *gin.Engine)
}
type ISystemService interface {
	Navigations(ctx context.Context) ([]*model.Navigation, error)
	PluginConfig(ctx context.Context) ([]pm3.PFrontend, error)
}
