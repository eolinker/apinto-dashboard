package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	groupUUID := ginCtx.Query("group")
	searchName := ginCtx.Query("search")

	pluginItems, err := p.modulePluginService.GetPlugins(ginCtx, groupUUID, searchName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugins fail. err:%s", err.Error())))
		return
	}

	pluginGroups, err := p.modulePluginService.GetPluginGroups(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugins fail. err:%s", err.Error())))
		return
	}
	plugins := make([]*dto.PluginListItem, 0, len(pluginItems))
	groups := make([]*dto.PluginGroup, 0, len(pluginGroups))

	for _, item := range pluginItems {
		plugins = append(plugins, &dto.PluginListItem{
			Id:      item.UUID,
			Name:    item.Name,
			Cname:   item.CName,
			Resume:  item.Resume,
			ICon:    item.ICon,
			Enable:  item.IsEnable,
			IsInner: item.IsInner,
		})
	}

	for _, group := range pluginGroups {
		groups = append(groups, &dto.PluginGroup{
			UUID: group.UUID,
			Name: group.Name,
		})
	}

	data := common.Map[string, interface{}]{}
	data["plugins"] = plugins
	data["groups"] = groups
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getPluginInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	pluginInfo, err := p.modulePluginService.GetPluginInfo(ginCtx, pluginUUID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugin Info fail. err:%s", err.Error())))
		return
	}

	data := common.Map[string, interface{}]{}
	info := &dto.PluginInfo{
		Id:        pluginInfo.UUID,
		Name:      pluginInfo.Name,
		Cname:     pluginInfo.CName,
		Resume:    pluginInfo.Resume,
		Icon:      pluginInfo.ICon,
		Enable:    pluginInfo.IsEnable,
		Uninstall: pluginInfo.Uninstall,
	}
	data["plugin"] = info
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getGroupsEnum(ginCtx *gin.Context) {
	pluginGroups, err := p.modulePluginService.GetPluginGroups(ginCtx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugin groups enum fail. err:%s", err.Error())))
		return
	}

	groups := make([]*dto.PluginGroup, 0, len(pluginGroups))
	for _, group := range pluginGroups {
		//排除内置插件目录
		if group.UUID == "inner_plugin" {
			continue
		}
		groups = append(groups, &dto.PluginGroup{
			UUID: group.UUID,
			Name: group.Name,
		})
	}
	data := common.Map[string, interface{}]{}
	data["groups"] = groups
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getEnableInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")
	info, err := p.modulePluginService.GetPluginEnableInfo(ginCtx, pluginUUID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error())))
		return
	}
	render, err := p.modulePluginService.GetPluginEnableRender(ginCtx, pluginUUID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error())))
		return
	}

	infoHeader := make([]dto.ExtendParams, 0, len(info.Header))
	infoQuery := make([]dto.ExtendParams, 0, len(info.Query))
	infoInitialize := make([]dto.ExtendParams, 0, len(info.Initialize))
	for _, h := range info.Header {
		infoHeader = append(infoHeader, dto.ExtendParams{
			Name:  h.Name,
			Value: h.Value,
		})
	}
	for _, q := range info.Query {
		infoQuery = append(infoQuery, dto.ExtendParams{
			Name:  q.Name,
			Value: q.Value,
		})
	}
	for _, i := range info.Initialize {
		infoInitialize = append(infoInitialize, dto.ExtendParams{
			Name:  i.Name,
			Value: i.Value,
		})
	}

	renderHeader := make([]dto.ExtendParamsRender, 0, len(render.Headers))
	renderQuery := make([]dto.ExtendParamsRender, 0, len(render.Querys))
	renderInitialize := make([]dto.ExtendParamsRender, 0, len(render.Initialize))
	for _, h := range render.Headers {
		renderHeader = append(renderHeader, dto.ExtendParamsRender{
			Name:        h.Name,
			Value:       h.Value,
			Title:       h.Title,
			Placeholder: h.Placeholder,
			Desc:        h.Desc,
		})
	}
	for _, q := range render.Querys {
		renderQuery = append(renderQuery, dto.ExtendParamsRender{
			Name:        q.Name,
			Value:       q.Value,
			Title:       q.Title,
			Placeholder: q.Placeholder,
			Desc:        q.Desc,
		})
	}
	for _, i := range render.Initialize {
		renderInitialize = append(renderInitialize, dto.ExtendParamsRender{
			Name:        i.Name,
			Value:       i.Value,
			Title:       i.Title,
			Placeholder: i.Placeholder,
			Desc:        i.Desc,
		})
	}

	enableInfo := &dto.PluginEnableInfo{
		Name:       info.Name,
		Navigation: info.Navigation,
		ApiGroup:   info.ApiGroup,
		Server:     info.Server,
		Header:     infoHeader,
		Query:      infoQuery,
		Initialize: infoInitialize,
	}

	enableRender := &dto.PluginEnableRender{
		Internet:   render.Internet,
		Invisible:  render.Invisible,
		ApiGroup:   render.ApiGroup,
		Headers:    renderHeader,
		Querys:     renderQuery,
		Initialize: renderInitialize,
	}

	data := common.Map[string, interface{}]{}
	data["module"] = enableInfo
	data["render"] = enableRender
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) install(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)
	pluginPackage, err := ginCtx.FormFile("plugin")
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("install plugin FormFilefail. err:%s", err.Error())))
		return
	}

	file, err := pluginPackage.Open()
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("install plugin openFile fail. err:%s", err.Error())))
		return
	}
	defer file.Close()

	// 检查文件类型和大小
	contentType := pluginPackage.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/gzip") {
		ginCtx.String(http.StatusBadRequest, "Invalid file type")
		return
	}
	if pluginPackage.Size > 4<<20 {
		ginCtx.String(http.StatusBadRequest, "File too large")
		return
	}

	groupName := ginCtx.PostForm("group_name")
	if groupName == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("install plugin fail. err: groupName is null. ")))
		return
	}

	fileBuffer := make([]byte, 0)
	_, err = file.Read(fileBuffer)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("install plugin read file fail. err:%s", err.Error())))
		return
	}

	//TODO 将压缩包存放本地

	err = p.modulePluginService.InstallPlugin(ginCtx, userId, groupName, fileBuffer)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("install plugin fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) enable(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	input := new(dto.PluginEnableInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	err := p.modulePluginService.EnablePlugin(ginCtx, pluginUUID, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Enable plugin fail. err:%s", err.Error())))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) disable(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	err := p.modulePluginService.DisablePlugin(ginCtx, pluginUUID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Disable plugin fail. err:%s", err.Error())))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
