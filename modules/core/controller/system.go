package controller

import (
	_ "embed"
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/custom"
	"github.com/eolinker/apinto-dashboard/initialize/inert"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/log"
	"gopkg.in/yaml.v3"
	"net/http"

	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"

	"github.com/eolinker/apinto-dashboard/modules/api"

	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/apinto-dashboard/controller"

	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/gin-gonic/gin"

	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

var (
	//go:embed frontend-plugin-inert.yml
	frontendPluginsInert []byte
)

func init() {
	var inertPluginFrontendConfig []pm3.PFrontend

	err := yaml.Unmarshal(frontendPluginsInert, &inertPluginFrontendConfig)
	if err != nil {
		panic(err)
	}
	inert.AddInertFrontendConfig(inertPluginFrontendConfig...)
}

type System struct {
	systemService  core.ISystemService
	clusterService cluster.IClusterService
	apiService     api.IAPIService
	dynamicService dynamic.IDynamicService
}

func systemRouters() apinto_module.RoutersInfo {
	s := &System{}

	bean.Autowired(&s.systemService)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.apiService)
	bean.Autowired(&s.dynamicService)
	return s.initRouter()
}

func (s *System) initRouter() apinto_module.RoutersInfo {
	routers := apinto_module.RoutersInfo{{
		Method:      http.MethodGet,
		Path:        "/api/system/modules",
		Authority:   pm3.Anonymous,
		HandlerFunc: s.list,
	},
		{
			Method:      http.MethodGet,
			Path:        "/api/system/quick_step",
			Authority:   pm3.Anonymous,
			HandlerFunc: s.quickStepInfo,
		}, {
			Method:      http.MethodGet,
			Path:        "/api/system/plugins",
			Authority:   pm3.Anonymous,
			HandlerFunc: s.fPlugins,
		},
	}
	return routers
}
func (s *System) fPlugins(ctx *gin.Context) {
	pluginConfigInstall, err := s.systemService.PluginConfig(ctx)
	if err != nil {
		log.Warn("get frontend plugin config error:", err)
	}
	inertPluginFrontendConfig := inert.GetFrontends()

	pluginConfig := make([]pm3.PFrontend, 0, len(pluginConfigInstall)+len(inertPluginFrontendConfig))
	pluginConfig = append(pluginConfig, inertPluginFrontendConfig...)
	if len(pluginConfigInstall) > 0 {
		pluginConfig = append(pluginConfig, pluginConfigInstall...)

	}
	data := make(map[string]any)
	data["plugins"] = pluginConfig
	for k, v := range version.BuildInfo() {
		data[k] = v
	}
	data["guide"] = custom.Guide()
	data["product"] = custom.Produce()
	data["powered"] = custom.Powered()
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(data))
	return
}
func (s *System) list(ctx *gin.Context) {
	list, err := s.systemService.Navigations(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"navigation": list,
		//"access":     access,
	}))
}

func (s *System) quickStepInfo(ctx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ctx)
	clusterCount, err := s.clusterService.ClusterCount(ctx, namespaceId)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	upstreamCount, err := s.dynamicService.Count(ctx, namespaceId, "service", nil)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	apiCount, err := s.apiService.APICount(ctx, namespaceId)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	apiPublishCount, err := s.apiService.APIOnlineCount(ctx, namespaceId)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"cluster":     clusterCount > 0,
		"upstream":    upstreamCount > 0,
		"api":         apiCount > 0,
		"publish_api": apiPublishCount > 0,
	}))
}
