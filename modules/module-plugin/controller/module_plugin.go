package controller

import (
	"bytes"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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

	pluginGroups, err := p.modulePluginService.GetPluginGroups()
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
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin Info fail. err:%s", err.Error()))
		return
	}

	data := common.Map[string, interface{}]{}
	info := &dto.PluginInfo{
		Id:        pluginInfo.UUID,
		Name:      pluginInfo.Name,
		Cname:     pluginInfo.CName,
		Resume:    pluginInfo.Resume,
		Icon:      pluginInfo.ICon,
		Enable:    pluginInfo.Enable,
		IsDisable: pluginInfo.IsDisable,
		Uninstall: pluginInfo.Uninstall,
	}
	data["plugin"] = info
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (p *modulePluginController) getGroupsEnum(ginCtx *gin.Context) {
	pluginGroups, err := p.modulePluginService.GetPluginGroups()
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin groups enum fail. err:%s", err.Error()))
		return
	}

	groups := make([]*dto.PluginGroup, 0, len(pluginGroups))
	for _, group := range pluginGroups {
		//排除内置插件目录
		//if group.UUID == "inner_plugin" {
		//	continue
		//}
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
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	render, err := p.modulePluginService.GetPluginEnableRender(ginCtx, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get plugin enable info fail. err:%s", err.Error()))
		return
	}

	//若模块名没有冲突，且没有需要填的，直接启用
	if !info.NameConflict && !render.Internet && len(render.Headers) == 0 && len(render.Querys) == 0 && len(render.Initialize) == 0 {
		userId := controller.GetUserId(ginCtx)
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
		Server:     info.Server,
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
			Placeholder: i.Placeholder,
			Desc:        i.Desc,
		})
	}

	enableRender := &dto.PluginEnableRender{
		Internet: render.Internet,
		//Invisible:  render.Invisible,
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
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin FormFilefail. err:%s", err.Error()))
		return
	}

	file, err := pluginPackage.Open()
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin openFile fail. err:%s", err.Error()))
		return
	}
	defer file.Close()

	// 检查文件类型和大小
	contentType := pluginPackage.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/x-gzip") {
		ginCtx.String(http.StatusBadRequest, "Invalid file type")
		return
	}
	if pluginPackage.Size > 4<<20 {
		ginCtx.String(http.StatusBadRequest, "File too large")
		return
	}

	groupName := ginCtx.PostForm("group_name")
	if groupName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin fail. err: groupName is null. "))
		return
	}

	//读取压缩文件的内容
	fileBuffer, err := io.ReadAll(file)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin read file fail. err:%s", err.Error()))
		return
	}
	packageFile := bytes.NewReader(fileBuffer)

	randomId := uuid.New()
	tmpDir := fmt.Sprintf("%s%s%s", PluginDir, string(os.PathSeparator), randomId)
	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		log.Error("安装插件失败, 无法创建目录:", err)
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin read file fail. err:%s", err.Error()))
		return
	}
	err = common.DeCompress(packageFile, tmpDir)
	if err != nil {
		//删除目录
		os.RemoveAll(tmpDir)
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin decompress file fail. err:%s", err.Error()))
		return
	}

	//校验解压目录下有没有必要的文件 plugin.yml icon README.md
	//读取plugin.yml
	pluginCfgFile, err := os.Open(path.Join(tmpDir, "plugin.yml"))
	if err != nil {
		//文件不存在，删除目录
		os.RemoveAll(tmpDir)
		controller.ErrorJson(ginCtx, http.StatusOK, "安装插件失败  配置文件文件不存在.")
		return
	}
	defer pluginCfgFile.Close()
	pluginBuffer, err := io.ReadAll(pluginCfgFile)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin read plugin.yml fail. err:%s", err.Error()))
		return
	}
	pluginCfg := new(model.PluginYmlCfg)
	err = yaml.Unmarshal(pluginBuffer, pluginCfg)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin unmarshal plugin.yml fail. err:%s", err.Error()))
		return
	}
	//TODO 校验plugin.yml

	//README.md是否存在
	mdFile, err := os.Open(path.Join(tmpDir, "README.md"))
	if err != nil {
		//文件不存在，删除目录
		os.RemoveAll(tmpDir)
		controller.ErrorJson(ginCtx, http.StatusOK, "安装插件失败 README.md 文件不存在")
		return
	}
	defer mdFile.Close()

	//图标文件是否存在
	iconFile, err := os.Open(path.Join(tmpDir, pluginCfg.ICon))
	if err != nil {
		//文件不存在，删除目录
		os.RemoveAll(tmpDir)
		controller.ErrorJson(ginCtx, http.StatusOK, "安装插件失败 图标文件不存在")
		return
	}
	defer iconFile.Close()

	err = p.modulePluginService.InstallPlugin(ginCtx, userId, pluginCfg, fileBuffer)
	if err != nil {
		//删除目录
		os.RemoveAll(tmpDir)
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("install plugin fail. err:%s", err.Error()))
		return
	}

	//将临时目录重命名为插件id
	newDirName := fmt.Sprintf("%s%s%s", PluginDir, string(os.PathSeparator), pluginCfg.ID)
	err = os.Rename(tmpDir, newDirName)
	if err != nil {
		log.Errorf("安装插件更改临时目录名失败 err: %s", err)
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) uninstall(ginCtx *gin.Context) {
	userId := controller.GetUserId(ginCtx)
	pluginUUID := ginCtx.Query("id")

	err := p.modulePluginService.UninstallPlugin(ginCtx, userId, pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	//删除插件在本地的缓存
	os.RemoveAll(fmt.Sprintf("%s%s%s", PluginDir, string(os.PathSeparator), pluginUUID))

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) enable(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	input := new(dto.PluginEnableInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	err := p.modulePluginService.EnablePlugin(ginCtx, controller.GetUserId(ginCtx), pluginUUID, input)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *modulePluginController) disable(ginCtx *gin.Context) {

	pluginUUID := ginCtx.Query("id")

	err := p.modulePluginService.DisablePlugin(ginCtx, controller.GetUserId(ginCtx), pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Disable plugin fail. err:%s", err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
