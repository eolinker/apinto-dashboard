package plugin_controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/plugin"
	plugin_dto "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type pluginClusterController struct {
	clusterPluginService plugin.IClusterPluginService
}

func newPluginClusterController() *pluginClusterController {

	p := &pluginClusterController{}
	bean.Autowired(&p.clusterPluginService)
	return p
}

// 插件列表
func (p *pluginClusterController) plugins(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)

	plugins, err := p.clusterPluginService.GetList(ginCtx, namespaceID, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get ClusterlPlugin List fail. err: %s", err.Error()))
		return
	}
	list := make([]*plugin_dto.CluPluginListItem, 0, len(plugins))
	for _, pluginInfo := range plugins {
		list = append(list, &plugin_dto.CluPluginListItem{
			Name:         pluginInfo.PluginName,
			Publish:      enum.PublishType(pluginInfo.Publish),
			Config:       pluginInfo.Config,
			Status:       enum.PluginStateType(pluginInfo.Status),
			ChangeStatus: enum.ChangeOptType(pluginInfo.ChangeState),
			ReleasedSort: pluginInfo.ReleasedSort,
			NowSort:      pluginInfo.NowSort,
			IsBuiltIn:    pluginInfo.IsBuiltIn,
			Operator:     pluginInfo.Operator,
			CreateTime:   common.TimeToStr(pluginInfo.CreateTime),
			UpdateTime:   common.TimeToStr(pluginInfo.UpdateTime),
		})
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{"plugins": list}))
}

// 获取插件配置
func (p *pluginClusterController) getPlugin(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	pluginName := ginCtx.Query("plugin_name")
	if pluginName == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprint("Get ClusterlPlugin Info fail. err: plugin_name can't be null. "))
		return
	}

	pluginInfo, err := p.clusterPluginService.GetPlugin(ginCtx, namespaceID, clusterName, pluginName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Get ClusterlPlugin Info fail. err: %s", err.Error()))
		return
	}

	info := &plugin_dto.ClusterPluginInfo{
		PluginName: pluginInfo.PluginName,
		Status:     enum.PluginStateType(pluginInfo.Status),
		Config:     pluginInfo.Config,
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{"info": info}))
}

// 配置插件
func (p *pluginClusterController) editPlugin(ginCtx *gin.Context) {
	clusterName := ginCtx.Param("cluster_name")
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)

	input := new(plugin_dto.ClusterPluginInfoInput)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	status, exist := enum.GetPluginState(input.Status)
	if !exist {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Edit ClusterPlugin fail. status %s is illegal. ", input.Status))
		return
	}

	err := p.clusterPluginService.EditPlugin(ginCtx, namespaceID, clusterName, userId, input.PluginName, status, input.Config)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Edit ClusterPlugin fail. err: %s", err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 变更历史列表
func (p *pluginClusterController) updateHistory(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	pageNumStr := ginCtx.Query("page_num")
	pageSizeStr := ginCtx.Query("page_size")
	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	histories, total, err := p.clusterPluginService.QueryHistory(ginCtx, namespaceId, pageNum, pageSize, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	list := make([]*plugin_dto.ClusterPluginHistoryItem, 0, len(histories))
	for _, h := range histories {
		pluginName := h.OldValue.PluginName
		if pluginName == "" {
			pluginName = h.NewValue.PluginName
		}
		list = append(list, &plugin_dto.ClusterPluginHistoryItem{
			Name: pluginName,
			OldConfig: plugin_dto.ClusterPluginConfig{
				Status: enum.PluginStateType(h.OldValue.Status),
				Config: h.OldValue.Config,
			},
			NewConfig: plugin_dto.ClusterPluginConfig{
				Status: enum.PluginStateType(h.NewValue.Status),
				Config: h.NewValue.Config,
			},
			CreateTime: common.TimeToStr(h.OptTime),
			OptType:    enum.ChangeOptType(h.OptType),
		})
	}

	m := common.Map[string, interface{}]{}
	m["histories"] = list
	m["total"] = total

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// 发布历史列表
func (p *pluginClusterController) publishHistory(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")
	pageNumStr := ginCtx.Query("page_num")
	pageSizeStr := ginCtx.Query("page_size")

	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	list, total, err := p.clusterPluginService.PublishHistory(ginCtx, namespaceId, pageNum, pageSize, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	histories := make([]*plugin_dto.ClusterPluginPublishItem, 0, len(list))
	for _, publish := range list {
		details := make([]*plugin_dto.ClusterPluginPublishDetails, 0, len(publish.Details))
		for _, detail := range publish.Details {
			details = append(details, &plugin_dto.ClusterPluginPublishDetails{
				Name: detail.Name,
				OldConfig: plugin_dto.ClusterPluginConfig{
					Status: enum.PluginStateType(detail.OldValue.Status),
					Config: detail.OldValue.Config,
				},
				NewConfig: plugin_dto.ClusterPluginConfig{
					Status: enum.PluginStateType(detail.NewValue.Status),
					Config: detail.NewValue.Config,
				},
				ReleasedSort: detail.ReleasedSort,
				NowSort:      detail.NowSort,
				OptType:      enum.ChangeOptType(detail.OptType),
				CreateTime:   common.TimeToStr(detail.CreateTime),
			})
		}
		histories = append(histories, &plugin_dto.ClusterPluginPublishItem{
			Id:         publish.Id,
			Name:       publish.Name,
			OptType:    enum.PublishOptType(publish.OptType),
			Operator:   publish.Operator,
			CreateTime: common.TimeToStr(publish.CreateTime),
			Details:    details,
		})
	}

	m := common.Map[string, interface{}]{}
	m["histories"] = histories
	m["total"] = total
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// 待发布的插件列表
func (p *pluginClusterController) toPublish(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	list, err := p.clusterPluginService.ToPublishes(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	bytes, _ := json.Marshal(list)
	source := common.Base64Encode(bytes)

	toPublishItems := make([]*plugin_dto.ClusterPluginToPublishItem, 0, len(list))
	isPublish := true
	for _, publish := range list {
		optType := enum.ChangeOptType(publish.OptType)
		toPublishItems = append(toPublishItems, &plugin_dto.ClusterPluginToPublishItem{
			Name:     publish.PluginName,
			Extended: publish.Extended,
			ReleasedConfig: plugin_dto.ClusterPluginConfig{
				Status: enum.PluginStateType(publish.ReleasedConfig.Status),
				Config: publish.ReleasedConfig.Config,
			},
			NoReleasedConfig: plugin_dto.ClusterPluginConfig{
				Status: enum.PluginStateType(publish.NoReleasedConfig.Status),
				Config: publish.NoReleasedConfig.Config,
			},
			ReleasedSort: publish.ReleasedSort,
			NowSort:      publish.NowSort,
			CreateTime:   common.TimeToStr(publish.CreateTime),
			OptType:      optType,
		})
	}

	plugins, err := p.clusterPluginService.GetList(ginCtx, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	m := common.Map[string, interface{}]{}

	defectPlugins := make([]string, 0)
	for _, pluginInfo := range plugins {
		if pluginInfo.Publish == 3 {
			defectPlugins = append(defectPlugins, pluginInfo.PluginName)
			isPublish = false
		}
	}

	if len(defectPlugins) > 0 {
		m["unpublished_msg"] = fmt.Sprintf("插件名为%s的插件处于缺失状态不可发布", strings.Join(defectPlugins, ","))
	}

	m["plugins"] = toPublishItems
	m["is_publish"] = isPublish
	m["source"] = source

	m["version_name"] = time.Now().Format("20060102150405") + "-release"

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}

// 发布插件
func (p *pluginClusterController) publish(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	clusterName := ginCtx.Param("cluster_name")

	input := &plugin_dto.ClusterPluginPublishInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	if input.VersionName == "" || input.Source == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "parameter error")
		return
	}
	background := ginCtx

	plugins, err := p.clusterPluginService.GetList(background, namespaceId, clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	defectPlugins := make([]string, 0)
	for _, pluginInfo := range plugins {
		if pluginInfo.Publish == 3 {
			defectPlugins = append(defectPlugins, pluginInfo.PluginName)
		}
	}

	if len(defectPlugins) > 0 {
		msg := fmt.Sprintf("插件名为%s的插件处于缺失状态不可发布", strings.Join(defectPlugins, ","))
		controller.ErrorJson(ginCtx, http.StatusOK, msg)
		return
	}

	userId := controller.GetUserId(ginCtx)
	if err = p.clusterPluginService.Publish(background, namespaceId, userId, clusterName, input.VersionName, input.Desc, input.Source); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
