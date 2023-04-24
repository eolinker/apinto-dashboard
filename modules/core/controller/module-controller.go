package controller

import (
	"github.com/eolinker/apinto-dashboard/frontend"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type ModuleController struct {
	coreService core.ICore
}

func NewModuleController() *ModuleController {
	m := &ModuleController{}
	bean.Autowired(&m.coreService)
	return m
}

func (m *ModuleController) HandleModule(ginCtx *gin.Context) {
	module := ginCtx.Param("module")
	path := ginCtx.Param("path")
	hasModule := m.coreService.HasModule(module, path)
	if hasModule {
		frontend.IndexHtml(ginCtx)
	} else {
		//ginCtx.Redirect(302, "/")
		frontend.IndexHtml(ginCtx)
	}

}
