package open_app_controller

import (
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	e := newExternalApplicationController()
	return apinto_module.RoutersInfo{{
		Method:      http.MethodGet,
		Path:        "/api/external-apps",
		Handler:     "applications.getList",
		HandlerFunc: []apinto_module.HandlerFunc{e.getList},
	},
		{
			Method:      http.MethodGet,
			Path:        "/api/external-app",
			Handler:     "applications.getInfo",
			HandlerFunc: []apinto_module.HandlerFunc{e.getInfo},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/external-app",
			Handler:     "applications.create",
			HandlerFunc: []apinto_module.HandlerFunc{e.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/external-app",
			Handler:     "applications.edit",
			HandlerFunc: []apinto_module.HandlerFunc{e.edit},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/external-app",
			Handler:     "applications.delete",
			HandlerFunc: []apinto_module.HandlerFunc{e.delete},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/external-app/enable",
			Handler:     "applications.enable",
			HandlerFunc: []apinto_module.HandlerFunc{e.enable},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/external-app/disable",
			Handler:     "applications.disable",
			HandlerFunc: []apinto_module.HandlerFunc{e.disable},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/external-app/token",
			Handler:     "applications.flushToken",
			HandlerFunc: []apinto_module.HandlerFunc{e.flushToken},
		},
	}
}
