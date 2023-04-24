package dynamic_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eolinker/eosc/log"

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

	modulePluginService module_plugin.IModulePlugin
	dynamicService      dynamic.IDynamicService
	clusterService      cluster.IClusterService

	Profession string
	Drivers    []*Basic
	Fields     []*Basic
	Skill      string
	Render     map[string]Render
}

func newDynamicController(name string, define interface{}) *dynamicController {
	tmp, _ := json.Marshal(define)
	var cfg DynamicDefine
	json.Unmarshal(tmp, &cfg)
	render := make(map[string]Render)
	for key, value := range cfg.Render {
		r := make(Render)
		err := json.Unmarshal([]byte(value), &r)
		if err != nil {
			log.Errorf("dynamic define parse error: %w,body is %s", err, value)
			continue
		}
		render[key] = r
	}

	d := &dynamicController{
		moduleName: name,
		Profession: cfg.Profession,
		Drivers:    cfg.Drivers,
		Fields:     cfg.Fields,
		Skill:      cfg.Skill,
		Render:     render,
	}
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
	pluginInfo, err := c.modulePluginService.GetEnabledPluginByModuleName(ctx, c.moduleName)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"id":      pluginInfo.UUID,
		"name":    pluginInfo.Name,
		"title":   pluginInfo.CName,
		"drivers": c.Drivers,
		"fields":  c.Fields,
		"list":    list,
	}))
	return
}

func (c *dynamicController) info(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	info, err := c.dynamicService.Info(ctx, namespaceID, c.Profession, uuid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(info))
}

func (c *dynamicController) online(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	clusters := ctx.GetStringSlice("cluster")
	userId := controller.GetUserId(ctx)
	success, fail, err := c.dynamicService.Online(ctx, namespaceID, c.Profession, uuid, clusters, userId)
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
	success, fail, err := c.dynamicService.Offline(ctx, namespaceID, c.Profession, uuid, clusters, userId)
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
	pluginInfo, err := c.modulePluginService.GetEnabledPluginByModuleName(ctx, c.moduleName)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"id":     pluginInfo.UUID,
		"name":   pluginInfo.Name,
		"title":  pluginInfo.CName,
		"render": c.Render,
	}))
}

func (c *dynamicController) clusterStatusList(ctx *gin.Context) {
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
	clusterInfo, err := c.dynamicService.ClusterStatuses(ctx, namespaceID, c.Profession, clusters, keyword, page, pageSize)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(clusterInfo))
}

func (c *dynamicController) clusterStatus(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uuid := ctx.Param("uuid")
	basic, clusters, err := c.dynamicService.ClusterStatus(ctx, namespaceID, c.Profession, uuid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"id":          basic.ID,
		"name":        basic.ID,
		"title":       basic.Title,
		"description": basic.Description,
		"clusters":    clusters,
	}))
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
	err = c.dynamicService.Create(ctx, namespaceID, c.Profession, worker.Title, worker.Name, worker.Driver, worker.Description, string(body), controller.GetUserId(ctx))
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
	err = c.dynamicService.Save(ctx, namespaceID, c.Profession, worker.Title, uuid, worker.Description, string(body), controller.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
