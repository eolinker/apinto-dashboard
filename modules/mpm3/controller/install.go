package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/dto"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/apinto-dashboard/pm3/pinstall"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// enablePluginStatus 当获取插件信息时可以直接启用时返回的状态码
const enablePluginStatus = 30001

type Install struct {
	installService mpm3.IInstallService
	pluginService  mpm3.IPluginService
}

func NewInstall() *Install {
	p := &Install{}
	bean.Autowired(&p.installService)
	bean.Autowired(&p.pluginService)
	return p
}
func (p *Install) apis() []pm3.Api {
	return []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/system/plugin/installed",
			Authority:   pm3.Private,
			HandlerFunc: p.pluginsInstalled,
		}, {
			Method:      http.MethodGet,
			Path:        "/api/system/plugin/info",
			Authority:   pm3.Private,
			HandlerFunc: p.getPluginInfo,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/system/plugin/groups/enum",
			Authority:   pm3.Private,
			HandlerFunc: p.getGroupsEnum,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/system/plugin/enable",
			Authority:   pm3.Private,
			HandlerFunc: p.getEnableInfo,
		}, {
			Method:      http.MethodPost,
			Path:        "/api/system/plugin/install",
			HandlerFunc: p.install,
		}, {
			Method:      http.MethodPost,
			Path:        "/api/system/plugin/uninstall",
			HandlerFunc: p.uninstall,
		},
	}
}

func (p *Install) install(ginCtx *gin.Context) {
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
	externPCfg, err := pinstall.Read(pluginYml)

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

	resources := model.NewPluginResources(externPCfg.ICon, "README.md", files)

	ctx := common.WidthUser(ginCtx, userId)
	err = p.installService.Install(ctx, externPCfg, resources)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("安装插件失败. err:%s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *Install) uninstall(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	err := p.installService.Uninstall(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *Install) pluginsInstalled(ginCtx *gin.Context) {
	groupUUID := ginCtx.Query("group")
	searchName := ginCtx.Query("search")

	pluginItems, err := p.pluginService.Search(ginCtx, groupUUID, searchName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugins fail. err:%s", err.Error()))
		return
	}

	pluginGroups, err := p.pluginService.GetGroups(ginCtx)
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
			Enable:  item.Enable,
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
func (p *Install) getPluginInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	pluginInfo, err := p.pluginService.GetPlugin(ginCtx, pluginUUID)
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

func (p *Install) getGroupsEnum(ginCtx *gin.Context) {
	pluginGroups, err := p.pluginService.GetGroups(ginCtx)
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

func (p *Install) getEnableInfo(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")
	info, err := p.pluginService.GetPluginEnableInfo(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	render, err := p.pluginService.GetEnableRender(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	//若模块名没有冲突，且没有需要填的，直接启用
	if len(render.Headers) == 0 && len(render.Querys) == 0 && len(render.Initialize) == 0 {
		userId := users.GetUserId(ginCtx)
		err = p.pluginService.EnablePlugin(ginCtx, userId, pluginUUID, &model.PluginEnableCfg{})
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("启用插件失败:%s", err.Error()))
			return
		}
		ginCtx.JSON(http.StatusOK, controller.NewResult(enablePluginStatus, nil, "success"))
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
	enableInfo := &dto.PluginEnableInfo{

		Header:     infoHeader,
		Query:      infoQuery,
		Initialize: infoInitialize,
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
		Headers: common.SliceToSlice(render.Headers, func(h model.ExtendParamsRender) dto.ExtendParamsRender {
			return dto.ExtendParamsRender{
				Name:        h.Name,
				Value:       h.Value,
				Title:       h.Title,
				Placeholder: h.Placeholder,
				Desc:        h.Desc,
			}
		}),
		Querys: common.SliceToSlice(render.Querys, func(h model.ExtendParamsRender) dto.ExtendParamsRender {
			return dto.ExtendParamsRender{
				Name:        h.Name,
				Value:       h.Value,
				Title:       h.Title,
				Placeholder: h.Placeholder,
				Desc:        h.Desc,
			}
		}),
		Initialize: common.SliceToSlice(render.Initialize, func(h model.ExtendParamsRender) dto.ExtendParamsRender {
			return dto.ExtendParamsRender{
				Name:        h.Name,
				Value:       h.Value,
				Title:       h.Title,
				Placeholder: h.Placeholder,
				Desc:        h.Desc,
			}
		}),
	}

	data := common.Map{}
	data["module"] = enableInfo
	data["render"] = enableRender
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
