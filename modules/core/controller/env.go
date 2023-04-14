package controller

import (
	env_controller "github.com/eolinker/apinto-dashboard/modules/base/env-controller"
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func envEnumRouters() apinto_module.RoutersInfo {
	ec := env_controller.NewEnumController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/enum/envs",
			Labels:      apinto_module.RouterLabelModule,
			HandlerFunc: []apinto_module.HandlerFunc{ec.GetEnv},
		},
	}
}
