package controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	env_controller "github.com/eolinker/apinto-dashboard/modules/base/env-controller"
	"net/http"
)

func envEnumRouters() apinto_module.RoutersInfo {
	ec := env_controller.NewEnumController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/enum/envs",
			HandlerFunc: []apinto_module.HandlerFunc{ec.GetEnv},
		},
	}
}
