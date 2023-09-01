package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
)

type modulePluginController struct {
	modulePluginService module_plugin.IModulePluginService
}

func newModulePluginController() *modulePluginController {

	p := &modulePluginController{}
	bean.Autowired(&p.modulePluginService)
	return p
}

// enablePluginStatus 当获取插件信息时可以直接启用时返回的状态码
const enablePluginStatus = 30001

// 插件列表
func (p *modulePluginController) plugins(ginCtx *gin.Context) {
	groupUUID := ginCtx.Query("group")
	searchName := ginCtx.Query("search")

	pluginItems, err := p.modulePluginService.GetPlugins(ginCtx, groupUUID, searchName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugins fail. err:%s", err.Error()))
		return
	}

	pluginGroups, err := p.modulePluginService.GetPluginGroups(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugins fail. err:%s", err.Error()))
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
			UUID:  group.UUID,
			Name:  group.Name,
			Count: group.Count,
		})
	}

	data := common.Map{}
	data["plugins"] = plugins
	data["groups"] = groups
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getPluginInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	pluginInfo, err := p.modulePluginService.GetPluginInfo(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin Info fail. err:%s", err.Error()))
		return
	}

	data := common.Map{}
	info := &dto.PluginInfo{
		Id:         pluginInfo.UUID,
		Name:       pluginInfo.Name,
		Cname:      pluginInfo.CName,
		Resume:     pluginInfo.Resume,
		Icon:       pluginInfo.ICon,
		Enable:     pluginInfo.Enable,
		CanDisable: pluginInfo.CanDisable,
		Uninstall:  pluginInfo.Uninstall,
	}
	data["plugin"] = info
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getGroupsEnum(ginCtx *gin.Context) {
	pluginGroups, err := p.modulePluginService.GetPluginGroups(ginCtx)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin groups enum fail. err:%s", err.Error()))
		return
	}

	groups := make([]*dto.PluginGroup, 0, len(pluginGroups))
	for _, group := range pluginGroups {
		groups = append(groups, &dto.PluginGroup{
			UUID: group.UUID,
			Name: group.Name,
		})
	}

	data := common.Map{}
	data["groups"] = groups
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getEnableInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")
	info, err := p.modulePluginService.GetPluginEnableInfo(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	render, err := p.modulePluginService.GetPluginEnableRender(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	//若模块名没有冲突，且没有需要填的，直接启用
	if !render.NameConflict && !render.Internet && len(render.Headers) == 0 && len(render.Querys) == 0 && len(render.Initialize) == 0 {
		userId := users.GetUserId(ginCtx)
		err = p.modulePluginService.EnablePlugin(ginCtx, userId, pluginUUID, &dto.PluginEnableInfo{
			Name:   info.Name,
			Server: "",
		})
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("启用插件失败:%s", err.Error()))
			return
		}
		ginCtx.JSON(http.StatusOK, controller.NewResult(enablePluginStatus, nil, "success"))
		return
	}

	//第一次启动时需要配置server, 则显示插件配置define里的server
	enableServer := info.Server
	if enableServer == "" && render.Server != "" && render.Internet {
		enableServer = render.Server
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
	enableInfo := &dto.PluginEnableInfo{
		Name:       info.Name,
		Server:     enableServer,
		Header:     infoHeader,
		Query:      infoQuery,
		Initialize: infoInitialize,
	}

	//若模块名没有冲突，不需要给前端发送
	if !render.NameConflict {
		enableInfo.Name = ""
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
			Type:        i.Type,
			Placeholder: i.Placeholder,
			Desc:        i.Desc,
		})
	}

	enableRender := &dto.PluginEnableRender{
		Internet:     render.Internet,
		NameConflict: render.NameConflict,
		//Invisible:  render.Invisible,
		Headers:    renderHeader,
		Querys:     renderQuery,
		Initialize: renderInitialize,
	}

	data := common.Map{}
	data["module"] = enableInfo
	data["render"] = enableRender
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) install(ginCtx *gin.Context) {
	userId := users.GetUserId(ginCtx)
	pluginPackage, err := ginCtx.FormFile("plugin")
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("从请求体读取插件安装包失败. err:%s", err.Error()))
		return
	}
	file, err := pluginPackage.Open()
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("插件安装包打开失败. err:%s", err.Error()))
		return
	}
	defer file.Close()

	// 检查文件类型和大小
	//contentType := pluginPackage.Header.Get("Content-Type")
	//if !strings.HasPrefix(contentType, "application/x-gzip") {
	//	ginCtx.String(http.StatusBadRequest, "Invalid file type")
	//	return
	//}
	//if pluginPackage.Size > 4<<20 {
	//	ginCtx.String(http.StatusBadRequest, "File too large")
	//	return
	//}

	//读取压缩文件的内容
	fileBuffer, err := io.ReadAll(file)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("读取插件安装包内容失败. err:%s", err.Error()))
		return
	}

	files, err := common.UnzipFromBytes(fileBuffer)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("解压插件安装包失败. err:%s", err.Error()))
		return
	}

	//校验解压目录下有没有必要的文件 plugin.yml icon README.md
	pluginYml, has := files["plugin.yml"]
	if !has {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("安装插件失败, plugin.yml不存在"))
		return
	}
	externPCfg := new(model.ExternPluginCfg)
	err = yaml.Unmarshal(pluginYml, externPCfg)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("plugin.yml解析失败. err:%s", err.Error()))
		return
	}
	//TODO 校验plugin.yml

	//README.md是否存在
	_, has = files["README.md"]
	if !has {
		controller.ErrorJson(ginCtx, http.StatusOK, "安装插件失败 README.md 文件不存在")
		return
	}

	//图标文件是否存在
	_, has = files[externPCfg.ICon]
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, "安装插件失败 图标文件不存在")
		return
	}
	pluginCfg := &model.PluginCfg{
		Version:    externPCfg.Version,
		Navigation: externPCfg.Navigation,
		GroupID:    externPCfg.GroupID,
		Resume:     externPCfg.Resume,
		Define:     externPCfg.Define,
	}
	resources := &model.PluginResources{
		Icon:   externPCfg.ICon,
		Readme: "README.md",
		Files:  files,
	}
	err = p.modulePluginService.InstallPlugin(ginCtx, userId, externPCfg.ID, externPCfg.Name, externPCfg.CName, externPCfg.Driver, externPCfg.ICon, pluginCfg, resources)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("安装插件失败. err:%s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) uninstall(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	err := p.modulePluginService.UninstallPlugin(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) enable(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	input := new(dto.PluginEnableInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	err := p.modulePluginService.EnablePlugin(ginCtx, users.GetUserId(ginCtx), pluginUUID, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) disable(ginCtx *gin.Context) {

	pluginUUID := ginCtx.Query("id")

	err := p.modulePluginService.DisablePlugin(ginCtx, users.GetUserId(ginCtx), pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Disable plugin fail. err:%s", err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
