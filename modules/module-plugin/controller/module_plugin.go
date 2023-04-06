package controller

import (
	"github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type modulePluginController struct {
	modulePluginService module_plugin.IModulePluginService
}

func RegisterModulePluginRouter(router gin.IRoutes) {

	p := &modulePluginController{}
	bean.Autowired(&p.modulePluginService)
	router.GET("/plugin/installed", p.plugins)
	router.GET("/plugin/info", p.getPluginInfo)
	router.GET("/plugin/groups/enum", p.getGroupsEnum)
	router.GET("/plugin/enable", p.getEnableInfo)

	router.POST("/plugin/install", p.install)
	router.POST("/plugin/enable", p.enable)
	router.POST("/plugin/disable", p.disable)

}

// 插件列表
func (p *modulePluginController) plugins(ginCtx *gin.Context) {

}

func (p *modulePluginController) getPluginInfo(ginCtx *gin.Context) {

}

func (p *modulePluginController) getGroupsEnum(ginCtx *gin.Context) {

}

func (p *modulePluginController) getEnableInfo(ginCtx *gin.Context) {

}

func (p *modulePluginController) install(ginCtx *gin.Context) {

}

func (p *modulePluginController) enable(ginCtx *gin.Context) {

}

func (p *modulePluginController) disable(ginCtx *gin.Context) {

}
