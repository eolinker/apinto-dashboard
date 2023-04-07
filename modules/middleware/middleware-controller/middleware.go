package middleware_controller

import (
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"net/http"

	"github.com/eolinker/apinto-dashboard/modules/middleware/dto"

	"github.com/eolinker/apinto-dashboard/enum"

	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/middleware"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type middlewareController struct {
	middlewareService   middleware.IMiddlewareService
	modulePluginService module_plugin.IModulePluginService
}

func RegisterMiddlewareGroupRouter(router gin.IRoutes) {
	c := &middlewareController{}
	bean.Autowired(&c.middlewareService)
	bean.Autowired(&c.modulePluginService)
	router.GET("/middleware", c.groups)
	router.POST("/middleware", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindMiddleware, c.save))
	router.PUT("/middleware", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindMiddleware, c.save))
}

func (m *middlewareController) groups(ginCtx *gin.Context) {
	groups, err := m.middlewareService.Groups(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	middlewares, err := m.modulePluginService.GetMiddlewareList(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())

		return
	}
	data := map[string]interface{}{
		"groups":      groups,
		"middlewares": middlewares,
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (m *middlewareController) save(ginCtx *gin.Context) {
	input := new(dto.MiddlewaresInput)

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if err := input.ValidCheck(); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := m.middlewareService.Save(ginCtx, input.Groups); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
