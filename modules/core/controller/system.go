package controller

import (
	"github.com/eolinker/apinto-dashboard/modules/core"
	"net/http"

	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"

	"github.com/eolinker/apinto-dashboard/modules/upstream"

	"github.com/eolinker/apinto-dashboard/modules/api"

	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/apinto-dashboard/controller"

	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/gin-gonic/gin"

	apinto_module "github.com/eolinker/apinto-module"
)

type System struct {
	navigationService core.INavigationService
	clusterService    cluster.IClusterService
	apiService        api.IAPIService
	upstreamService   upstream.IService

	routers apinto_module.RoutersInfo
}

func newSystem() *System {
	s := &System{routers: make(apinto_module.RoutersInfo, 0)}
	s.initRouter()
	bean.Autowired(&s.navigationService)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.apiService)
	bean.Autowired(&s.upstreamService)
	return s
}

func (s *System) RoutersInfo() apinto_module.RoutersInfo {
	return s.routers
}

func (s *System) initRouter() {
	s.routers = append(s.routers, apinto_module.RouterInfo{
		Method:      http.MethodGet,
		Path:        "/api/system/modules",
		Handler:     "core.modules",
		Labels:      apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{s.list},
	}, apinto_module.RouterInfo{
		Method:      http.MethodGet,
		Path:        "/api/system/quick_step",
		Handler:     "core.quick_step",
		Labels:      apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{s.quickStepInfo},
	})
}

func (s *System) list(ctx *gin.Context) {
	list, access, err := s.navigationService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"navigation": list,
		"access":     access,
	}))
}

func (s *System) quickStepInfo(ctx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ctx)
	clusterCount, err := s.clusterService.ClusterCount(ctx, namespaceId)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}
	upstreamCount, err := s.upstreamService.UpstreamCount(ctx, namespaceId)
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
