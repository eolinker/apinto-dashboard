package navigation_controller

import (
	"net/http"

	navigation_dto "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-dto"

	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/navigation"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type navigationController struct {
	navigationService navigation.INavigationService
}

func newNavigationController() *navigationController {
	c := &navigationController{}
	bean.Autowired(&c.navigationService)
	return c
}

//func RegisterNavigationRouter(router gin.IRoutes) {
//	c := &navigationController{}
//	bean.Autowired(&c.navigationService)
//	router.GET("/system/navigation", c.list)
//	router.GET("/system/navigation/:uuid", c.info)
//	router.POST("/system/navigation", controller.AuditLogHandler(enum.LogOperateTypeCreate, enum.LogKindNavigation, c.add))
//	router.PUT("/system/navigation/:uuid", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindNavigation, c.update))
//	router.DELETE("/system/navigation/:uuid", controller.AuditLogHandler(enum.LogOperateTypeDelete, enum.LogKindCommonGroup, c.delete))
//	router.PUT("/system/navigation", controller.AuditLogHandler(enum.LogOperateTypeEdit, enum.LogKindNavigation, c.sort))
//}

func (n *navigationController) list(ctx *gin.Context) {
	navigations, err := n.navigationService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"navigations": navigations,
	}))
}

func (n *navigationController) info(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	navigation, err := n.navigationService.Info(ctx, uuid)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"navigation": navigation,
	}))
}

func (n *navigationController) add(ctx *gin.Context) {
	input := new(navigation_dto.Navigation)

	if err := ctx.BindJSON(input); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if err := n.navigationService.Add(ctx, input.Uuid, input.Name, input.Icon); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (n *navigationController) update(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	input := new(navigation_dto.Navigation)

	if err := ctx.BindJSON(input); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if err := n.navigationService.Save(ctx, uuid, input.Name, input.Icon); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (n *navigationController) delete(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if err := n.navigationService.Delete(ctx, uuid); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (n *navigationController) sort(ctx *gin.Context) {
	uuids := make([]string, 0)
	if err := ctx.BindJSON(&uuids); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	if err := n.navigationService.Sort(ctx, uuids); err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
