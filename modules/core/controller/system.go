package controller

import (
	"net/http"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/modules/navigation"

	"github.com/gin-gonic/gin"

	apinto_module "github.com/eolinker/apinto-module"
)

type System struct {
	navigationService   navigation.INavigationService
	modulePluginService module_plugin.IModulePlugin
	routers             apinto_module.RoutersInfo
}

func (s *System) RoutersInfo() apinto_module.RoutersInfo {
	return s.routers
}

func (s *System) initRouter() {
	if s.routers == nil {
		s.routers = make(apinto_module.RoutersInfo, 0)
	}
	s.routers = append(s.routers, apinto_module.RouterInfo{
		Method:      http.MethodGet,
		Path:        "/api/system/modules",
		Handler:     "core.modules",
		Labels:      apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{s.modules},
	})
}

func (s *System) modules(ctx *gin.Context) {
	list, err := s.navigationService.List(ctx)
	if err != nil {

	}
}
