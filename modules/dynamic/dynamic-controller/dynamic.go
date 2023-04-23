package dynamic_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	dynamic_dto "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-dto"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/controller"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

type dynamicController struct {
	moduleName string

	modulePluginService module_plugin.IModulePluginService
	dynamicService      dynamic.IDynamicService
	clusterService      cluster.IClusterService
	*DynamicModulePlugin
}

func newDynamicController(name string, define *DynamicModulePlugin) *dynamicController {
	d := &dynamicController{moduleName: name, DynamicModulePlugin: define}
	bean.Autowired(&d.dynamicService)
	return d
}

func (c *dynamicController) list(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	keyword := ctx.GetString("keyword")
	clusters := ctx.GetStringSlice("cluster")
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("page_size")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		page = 15
	}

	cs, err := c.clusterService.GetByNames(ctx, namespaceID, clusters)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}

	columns := make([]string, 0, len(c.Fields))
	fields := make([]*Basic, 0, len(c.Fields)+len(clusters)+len(defaultFields))
	for _, field := range c.Fields {
		columns = append(columns, field.Name)
		fields = append(fields, field)
	}
	for _, cc := range cs {
		fields = append(fields, &Basic{
			Name:  cc.Name,
			Title: fmt.Sprintf("状态：%s", cc.Name),
			Attr:  "status",
			Enum: []string{
				"已发布",
				"待发布",
				"未发布",
			},
		})
	}
	fields = append(fields, defaultFields...)

	list, err := c.dynamicService.List(ctx, namespaceID, c.moduleName, columns, keyword, page, pageSize)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"id":      "",
		"name":    "",
		"title":   "",
		"drivers": c.Drivers,
		"fields":  c.Fields,
		"list":    list,
	}))
	return
}

func (c *dynamicController) info(ctx *gin.Context) {

}

func (c *dynamicController) online(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	clusters := ctx.GetStringSlice("cluster")
	userId := controller.GetUserId(ctx)
	success, fail, err := c.dynamicService.Online(ctx, namespaceID, c.moduleName, c.Profession, uuid, clusters, userId)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	if len(fail) > 0 {
		ctx.JSON(http.StatusOK, controller.NewResult(-1, map[string]interface{}{
			"success": success,
			"fail":    fail,
		}, "online error"))
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *dynamicController) offline(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	clusters := ctx.GetStringSlice("cluster")
	userId := controller.GetUserId(ctx)
	success, fail, err := c.dynamicService.Offline(ctx, namespaceID, c.moduleName, c.Profession, uuid, clusters, userId)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	if len(fail) > 0 {
		ctx.JSON(http.StatusOK, controller.NewResult(-1, map[string]interface{}{
			"success": success,
			"fail":    fail,
		}, "offline error"))
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *dynamicController) render(ctx *gin.Context) {
}

func (c *dynamicController) clusterStatusList(ctx *gin.Context) {
}

func (c *dynamicController) clusterStatus(ctx *gin.Context) {
}

func (c *dynamicController) batchDelete(ctx *gin.Context) {

}

func (c *dynamicController) create(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	var worker dynamic_dto.WorkerInfo
	err := ctx.BindJSON(&worker)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	body, _ := json.Marshal(worker.Append)
	err = c.dynamicService.Create(ctx, namespaceID, c.moduleName, worker.Title, worker.Name, worker.Driver, worker.Description, string(body), controller.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *dynamicController) save(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	var worker dynamic_dto.WorkerInfo
	err := ctx.BindJSON(&worker)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	body, _ := json.Marshal(worker.Append)
	err = c.dynamicService.Save(ctx, namespaceID, c.moduleName, worker.Title, uuid, worker.Description, string(body), controller.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
