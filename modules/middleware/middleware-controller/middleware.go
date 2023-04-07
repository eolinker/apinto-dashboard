package middleware_controller

import (
	"net/http"

	"github.com/eolinker/apinto-dashboard/modules/middleware/dto"

	"github.com/eolinker/apinto-dashboard/enum"

	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/middleware"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type middlewareController struct {
	middlewareService middleware.IMiddlewareService
}

func RegisterMiddlewareGroupRouter(router gin.IRoutes) {
	c := &middlewareController{}
	bean.Autowired(&c.middlewareService)
	router.GET("/system/middleware", c.groups)
	router.POST("/system/middleware", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindMiddleware), c.save)
	router.PUT("/system/middleware", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindMiddleware), c.save)
}

func (m *middlewareController) groups(ginCtx *gin.Context) {
	groups, err := m.middlewareService.Groups(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(groups))
}

func (m *middlewareController) save(ginCtx *gin.Context) {
	input := new(dto.Middlewares)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	if err := input.ValidCheck(); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if err := m.middlewareService.Save(ginCtx, input.String()); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
