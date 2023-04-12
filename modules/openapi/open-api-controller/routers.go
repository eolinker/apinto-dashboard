package open_api_controller

import (
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	c := newOpenApiController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/apis/import",
			Handler:     "applications.getImportInfo",
			HandlerFunc: []apinto_module.HandlerFunc{c.getImportInfo},
			Labels:      apinto_module.RouterLabelOpenApi,
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/apis/import",
			Handler:     "applications.syncAPI",
			HandlerFunc: []apinto_module.HandlerFunc{c.syncAPI},
			Labels:      apinto_module.RouterLabelOpenApi,
		},
	}
}
