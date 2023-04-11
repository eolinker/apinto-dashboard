package controller

import (
	"net/http"

	"github.com/eolinker/apinto-dashboard/controller"

	"github.com/eolinker/apinto-dashboard/modules/core/service"

	"github.com/gin-gonic/gin"

	apinto_module "github.com/eolinker/apinto-module"
)

type System struct {
	navigationService service.INavigationService
	routers           apinto_module.RoutersInfo
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
		HandlerFunc: []apinto_module.HandlerFunc{s.list},
	})
}

func (s *System) list(ctx *gin.Context) {
	list, err := s.navigationService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.NewErrorResult(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"navigation": list,
		"access":     "",
	}))
}
