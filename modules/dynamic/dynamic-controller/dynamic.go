package dynamic_controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"net/http"
	"strconv"

	cluster_model "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"

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
	drivers             []string

	Profession string
	Drivers    []*Basic
	Fields     []*Basic
	Skill      string
	Render     map[string]Render
}

func newDynamicController(name string, define interface{}) *dynamicController {
	//tmp, _ := json.Marshal(define)
	var cfg DynamicDefine
	json.Unmarshal(define.([]byte), &cfg)
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
	drivers := make([]string, 0, len(cfg.Drivers))
	for _, driver := range cfg.Drivers {
		drivers = append(drivers, driver.Name)
	}
	d := &dynamicController{
		moduleName: name,
		drivers:    drivers,
		Profession: cfg.Profession,
		Drivers:    cfg.Drivers,
		Fields:     cfg.Fields,
		Skill:      cfg.Skill,
		Render:     render,
	}
	bean.Autowired(&d.dynamicService)
	bean.Autowired(&d.modulePluginService)
	bean.Autowired(&d.clusterService)
	return d
}

func (c *dynamicController) getPage(ctx *gin.Context) (int, int) {
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pz, _ := strconv.Atoi(pageSize)
	if pz < 1 {
		pz = 15
	}
	return p, pz
}

func (c *dynamicController) list(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	keyword := ctx.Query("keyword")
	clusterNames := ctx.Query("cluster")
	names := make([]string, 0)
	if clusterNames != "" {
		err := json.Unmarshal([]byte(clusterNames), &names)
		if err != nil {
			controller.ErrorJson(ctx, http.StatusOK, err.Error())
			return
		}
	}
	page, pageSize := c.getPage(ctx)
	all := len(names) < 1
	var err error
	var cs []*cluster_model.Cluster
	if all {
		cs, err = c.clusterService.GetAllCluster(ctx)
	} else {
		cs, err = c.clusterService.GetByNames(ctx, namespaceID, names)
	}
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}

	columns := make([]string, 0, len(c.Fields))
	fields := make([]*Basic, 0, len(c.Fields)+len(names)+len(defaultFields))
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

	list, total, err := c.dynamicService.List(ctx, namespaceID, c.Profession, columns, c.drivers, keyword, page, pageSize)
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
		"fields":  fields,
		"list":    list,
		"total":   total,
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
	var tmp dynamic_dto.Cluster
	err := ctx.BindJSON(&tmp)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	userId := users.GetUserId(ctx)
	success, fail, err := c.dynamicService.Online(ctx, namespaceID, c.Profession, uuid, tmp.Cluster, userId)
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
	var tmp dynamic_dto.Cluster
	err := ctx.BindJSON(&tmp)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	userId := users.GetUserId(ctx)
	success, fail, err := c.dynamicService.Offline(ctx, namespaceID, c.Profession, uuid, tmp.Cluster, userId)
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
	keyword := ctx.Query("keyword")
	clusterNames := ctx.Query("cluster")
	names := make([]string, 0)
	if clusterNames != "" {
		err := json.Unmarshal([]byte(clusterNames), &names)
		if err != nil {
			controller.ErrorJson(ctx, http.StatusOK, err.Error())
			return
		}
	}

	page, pageSize := c.getPage(ctx)
	clusterInfo, err := c.dynamicService.ClusterStatuses(ctx, namespaceID, c.Profession, names, c.drivers, keyword, page, pageSize)
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
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	ids := ctx.Query("uuids")
	uuids := make([]string, 0)
	err := json.Unmarshal([]byte(ids), &uuids)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	success := make([]string, 0, len(uuids))
	fail := make([]string, 0, len(uuids))
	for _, uuid := range uuids {
		info, _, err := c.dynamicService.ClusterStatus(ctx, namespaceID, c.Profession, uuid)
		if err != nil {
			fail = append(fail, uuid)
		}
		if info.Online {
			fail = append(fail, uuid)
			continue
		}
		err = c.dynamicService.Delete(ctx, namespaceID, c.Profession, uuid)
		if err != nil {
			fail = append(fail, uuid)
		} else {
			success = append(success, uuid)
		}
	}
	if len(fail) > 0 {
		ctx.JSON(http.StatusOK, controller.NewResult(-1, map[string]interface{}{
			"success": success,
			"fail":    fail,
		}, "delete error"))
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
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
	err = c.dynamicService.Create(ctx, namespaceID, c.Profession, worker.Title, worker.Id, worker.Driver, worker.Description, string(body), users.GetUserId(ctx))
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
	err = c.dynamicService.Save(ctx, namespaceID, c.Profession, worker.Title, uuid, worker.Description, string(body), users.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
