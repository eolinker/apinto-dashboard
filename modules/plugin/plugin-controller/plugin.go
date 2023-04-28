package plugin_controller

import (
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	plugin2 "github.com/eolinker/apinto-dashboard/client/v1/initialize/plugin"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/plugin"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-dto"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	plugin_service "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pluginController struct {
	pluginService plugin.IPluginService
	extenderCache plugin_service.IExtenderCache
}

func newPluginController() *pluginController {
	p := &pluginController{}
	bean.Autowired(&p.pluginService)
	bean.Autowired(&p.extenderCache)
	return p
}

// 单个插件信息
func (p *pluginController) plugin(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	name := ginCtx.Query("name")

	pluginInfo, err := p.pluginService.GetByName(ginCtx, namespaceId, name)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	data := common.Map[string, interface{}]{}
	data["plugin"] = plugin_dto.PluginItem{
		Name:     pluginInfo.Name,
		Extended: pluginInfo.Extended,
		Desc:     pluginInfo.Desc,
		Rely:     pluginInfo.RelyName,
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 插件列表
func (p *pluginController) plugins(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := p.pluginService.GetList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resultList := make([]plugin_dto.PluginListItem, 0, len(list))

	for _, pluginInfo := range list {
		updateTime := ""
		if !pluginInfo.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(pluginInfo.UpdateTime)
		}

		resultList = append(resultList, plugin_dto.PluginListItem{
			Name:       pluginInfo.Name,
			Extended:   pluginInfo.Extended,
			Desc:       pluginInfo.Desc,
			UpdateTime: updateTime,
			Operator:   pluginInfo.OperatorStr,
			IsDelete:   pluginInfo.IsDelete,
			IsBuilt:    pluginInfo.IsBuilt,
		})
	}

	data := common.Map[string, interface{}]{}
	data["plugins"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 只有基本信息的插件列表
func (p *pluginController) basicInfoPlugins(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	list, err := p.pluginService.GetBasicInfoList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resultList := make([]plugin_dto.PluginListItem, 0, len(list))

	for _, pluginInfo := range list {
		resultList = append(resultList, plugin_dto.PluginListItem{
			Name:     pluginInfo.Name,
			Extended: pluginInfo.Extended,
			Desc:     pluginInfo.Desc,
		})
	}

	data := common.Map[string, interface{}]{}
	data["plugins"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 新增插件
func (p *pluginController) createPlugin(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	input := new(plugin_dto.PluginInput)

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//校验名称是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	pluginInfo := &plugin_model.PluginInput{
		Name:     input.Name,
		Extended: input.Extended,
		RelyName: input.Rely,
		Desc:     input.Desc,
	}

	if err := p.pluginService.Create(ginCtx, namespaceId, userId, pluginInfo); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 修改插件
func (p *pluginController) updatePlugin(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	input := new(plugin_dto.PluginInput)

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//校验名称是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	pluginInfo := &plugin_model.PluginInput{
		Name:     input.Name,
		Extended: input.Extended,
		RelyName: input.Rely,
		Desc:     input.Desc,
	}

	if err := p.pluginService.Update(ginCtx, namespaceId, userId, pluginInfo); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 删除插件
func (p *pluginController) delPlugin(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	name := ginCtx.Query("name")

	if err := p.pluginService.Delete(ginCtx, namespaceId, userId, name); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 插件的扩展ID列表
func (p *pluginController) pluginExtendeds(ginCtx *gin.Context) {
	extenderList, err := p.extenderCache.GetAll(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	//剔除黑名单中的扩展ID
	blackExtendedPluginConf := plugin2.GetBlackExtendedPluginConf()
	blackMaps := common.SliceToMap(blackExtendedPluginConf, func(t *v1.GlobalPlugin) string {
		return t.Id
	})

	names := make([]string, 0)
	for _, extender := range extenderList {
		if _, ok := blackMaps[extender.Id]; ok {
			continue
		}
		names = append(names, extender.Id)
	}

	data := common.Map[string, interface{}]{}
	data["extendeds"] = names
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 修改插件顺序
func (p *pluginController) pluginSort(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	pluginSort := new(plugin_dto.PluginSort)

	if err := ginCtx.BindJSON(pluginSort); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := p.pluginService.Sort(ginCtx, namespaceId, userId, pluginSort.Names); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 获取作为选项的插件列表（剔除掉黑名单）
func (p *pluginController) pluginEnum(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	plugins, err := p.pluginService.GetList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//剔除黑名单中的扩展ID
	blackExtendedPluginConf := plugin2.GetBlackExtendedPluginConf()
	blackMaps := common.SliceToMap(blackExtendedPluginConf, func(t *v1.GlobalPlugin) string {
		return t.Id
	})

	resultList := make([]plugin_dto.PluginEnum, 0)
	for _, pluginInfo := range plugins {
		if _, ok := blackMaps[pluginInfo.Extended]; ok {
			continue
		}
		resultList = append(resultList, plugin_dto.PluginEnum{
			Name:   pluginInfo.Name,
			Config: pluginInfo.Schema,
		})

	}
	data := common.Map[string, interface{}]{}
	data["plugins"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 插件配置信息
func (p *pluginController) pluginRender(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	plugins, err := p.pluginService.GetList(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resultList := make([]plugin_dto.PluginListItem, 0, len(plugins))

	for _, pluginInfo := range plugins {
		resultList = append(resultList, plugin_dto.PluginListItem{
			Name:   pluginInfo.Name,
			Config: pluginInfo.Schema,
		})
	}

	data := common.Map[string, interface{}]{}
	data["plugins"] = resultList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
