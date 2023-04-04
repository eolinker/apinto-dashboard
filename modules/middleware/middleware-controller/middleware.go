package middleware_controller

import (
	"net/http"

	"github.com/eolinker/apinto-dashboard/modules/middleware/dto"

	"github.com/eolinker/apinto-dashboard/enum"

	"github.com/eolinker/apinto-dashboard/controller"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
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
	router.GET("/middleware/group", c.groups)
	router.GET("/middleware/group/:uuid", c.info)
	router.POST("/middleware/group", controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindMiddlewareGroup), c.add)
	router.PUT("/middleware/group/:uuid", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindMiddlewareGroup), c.update)
	router.DELETE("/middleware/group/:uuid", controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindMiddlewareGroup), c.delete)
}

func (m *middlewareController) groups(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groups, err := m.middlewareService.GroupList(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(groups))
}

func (m *middlewareController) info(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Param("uuid")
	info, err := m.middlewareService.GroupInfo(ginCtx, namespaceId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(info))
}

func (m *middlewareController) add(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	input := new(dto.MiddlewareGroup)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	operator := controller.GetUserId(ginCtx)
	err := m.middlewareService.CreateGroup(ginCtx, namespaceId, operator, input.ID, input.Prefix, input.Middlewares)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (m *middlewareController) update(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Param("uuid")
	input := new(dto.MiddlewareGroup)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	operator := controller.GetUserId(ginCtx)
	err := m.middlewareService.UpdateGroup(ginCtx, namespaceId, operator, uuid, input.Prefix, input.Middlewares)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (m *middlewareController) delete(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Param("uuid")
	operator := controller.GetUserId(ginCtx)
	err := m.middlewareService.DeleteGroup(ginCtx, namespaceId, operator, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
