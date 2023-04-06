package plugin_template_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/online/online-dto"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-dto"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pluginTemplateController struct {
	pluginTemplateService plugin_template.IPluginTemplateService
}

func RegisterPluginTemplateRouter(router gin.IRoutes) {
	p := &pluginTemplateController{}
	bean.Autowired(&p.pluginTemplateService)

	router.GET("/plugin/templates", p.templates)
	router.GET("/plugin/template/enum", p.templateEnum)

	router.POST("/plugin/template", controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindPluginTemplate), p.createTemplate)
	router.PUT("/plugin/template", controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindPluginTemplate), p.updateTemplate)
	router.DELETE("/plugin/template", controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindPluginTemplate), p.delTemplate)

	router.GET("/plugin/template", p.template)

	router.GET("/plugin/template/onlines", p.onlines)
	router.PUT("/plugin/template/online", controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindPluginTemplate), p.online)
	router.PUT("/plugin/template/offline", controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindPluginTemplate), p.offline)
}

// 插件模板列表
func (p *pluginTemplateController) templates(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	templates, err := p.pluginTemplateService.GetList(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	resultList := make([]*plugin_template_dto.PluginTemplate, 0, len(templates))

	for _, template := range templates {
		createTime := ""
		updateTime := ""
		if !template.CreateTime.IsZero() {
			createTime = common.TimeToStr(template.CreateTime)
		}
		if !template.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(template.UpdateTime)
		}
		pluginTemplate := &plugin_template_dto.PluginTemplate{
			UUID:       template.UUID,
			Name:       template.Name,
			Desc:       template.Desc,
			Operator:   template.OperatorStr,
			CreateTime: createTime,
			UpdateTime: updateTime,
			IsDelete:   template.IsDelete,
		}
		resultList = append(resultList, pluginTemplate)
	}

	data := common.Map[string, interface{}]{}
	data["templates"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *pluginTemplateController) templateEnum(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	templates, err := p.pluginTemplateService.GetUsableList(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	resultList := make([]*plugin_template_dto.PluginTemplate, 0, len(templates))

	for _, template := range templates {
		pluginTemplate := &plugin_template_dto.PluginTemplate{
			UUID: template.UUID,
			Name: template.Name,
		}
		resultList = append(resultList, pluginTemplate)
	}

	data := common.Map[string, interface{}]{}
	data["templates"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 新增插件模板
func (p *pluginTemplateController) createTemplate(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	input := new(plugin_template_dto.PluginTemplateInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if input.Name == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("插件模板名称为必填项"))
		return
	}

	if len(input.Plugins) == 0 {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("插件配置为必填项"))
		return
	}

	//校验名称是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	plugins := make([]*plugin_template_model.PluginInfo, 0)
	for _, plugin := range input.Plugins {
		plugins = append(plugins, &plugin_template_model.PluginInfo{
			Name:    plugin.Name,
			Config:  plugin.Config,
			Disable: plugin.Disable,
		})
	}

	detail := &plugin_template_model.PluginTemplateDetail{
		PluginTemplate: &plugin_template_model.PluginTemplate{
			PluginTemplate: &plugin_template_entry.PluginTemplate{
				Name: input.Name,
				Desc: input.Desc,
			},
		},
		Plugins: plugins,
	}

	if err := p.pluginTemplateService.Create(ginCtx, namespaceId, userId, detail); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 修改插件模板
func (p *pluginTemplateController) updateTemplate(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	input := new(plugin_template_dto.PluginTemplateInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if len(input.Plugins) == 0 {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult("插件配置为必填项"))
		return
	}

	//校验名称是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	plugins := make([]*plugin_template_model.PluginInfo, 0)
	for _, plugin := range input.Plugins {
		plugins = append(plugins, &plugin_template_model.PluginInfo{
			Name:    plugin.Name,
			Config:  plugin.Config,
			Disable: plugin.Disable,
		})
	}

	detail := &plugin_template_model.PluginTemplateDetail{
		PluginTemplate: &plugin_template_model.PluginTemplate{
			PluginTemplate: &plugin_template_entry.PluginTemplate{
				UUID: input.Uuid,
				Desc: input.Desc,
			},
		},
		Plugins: plugins,
	}

	if err := p.pluginTemplateService.Update(ginCtx, namespaceId, userId, detail); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 删除插件模板
func (p *pluginTemplateController) delTemplate(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	uuid := ginCtx.Query("uuid")

	if err := p.pluginTemplateService.Delete(ginCtx, namespaceId, userId, uuid); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 单个插件模板信息
func (p *pluginTemplateController) template(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	uuid := ginCtx.Query("uuid")
	pluginTemplateDetail, err := p.pluginTemplateService.GetByUUID(ginCtx, namespaceId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	plugins := make([]*plugin_template_dto.PluginInfo, 0, len(pluginTemplateDetail.Plugins))
	for _, plugin := range pluginTemplateDetail.Plugins {
		plugins = append(plugins, &plugin_template_dto.PluginInfo{
			Name:    plugin.Name,
			Config:  plugin.Config,
			Disable: plugin.Disable,
		})
	}
	pluginTemplateOutput := &plugin_template_dto.PluginTemplateOutput{
		Name:    pluginTemplateDetail.Name,
		Desc:    pluginTemplateDetail.Desc,
		Plugins: plugins,
	}

	data := common.Map[string, interface{}]{}
	data["template"] = pluginTemplateOutput

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))

}

// 上线管理列表
func (p *pluginTemplateController) onlines(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	list, err := p.pluginTemplateService.OnlineList(ginCtx, namespaceId, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	resp := make([]*online_dto.OnlineOut, 0, len(list))
	for _, item := range list {
		updateTime := ""
		if !item.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(item.UpdateTime)
		}
		resp = append(resp, &online_dto.OnlineOut{
			Name:       item.ClusterName,
			Env:        item.ClusterEnv,
			Status:     enum.OnlineStatus(item.Status),
			Disable:    item.Disable,
			Operator:   item.Operator,
			UpdateTime: updateTime,
		})
	}

	m := make(map[string]interface{})
	m["clusters"] = resp
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// 上线管理-上线/更新
func (p *pluginTemplateController) online(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("uuid can't be nil")))
		return
	}
	input := &online_dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	router, err := p.pluginTemplateService.Online(ginCtx, namespaceId, userId, uuid, input.ClusterName)
	if err != nil && router == nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	} else if err == nil {
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
		return
	}

	msg := ""
	if err != nil {
		msg = err.Error()
	}

	m := make(map[string]interface{})
	m["router"] = router
	ginCtx.JSON(http.StatusOK, controller.Result{
		Code: -1,
		Data: m,
		Msg:  msg,
	})

}

// 下线
func (p *pluginTemplateController) offline(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("uuid can't be nil")))
		return
	}
	input := &online_dto.UpdateOnlineStatusInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	if err := p.pluginTemplateService.Offline(ginCtx, namespaceId, userId, uuid, input.ClusterName); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
