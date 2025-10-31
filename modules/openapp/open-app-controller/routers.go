package open_app_controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/apinto-dashboard/pm3/middleware"
	"net/http"
	"strings"
)

func initRouter(name string) (apinto_module.RoutersInfo, []apinto_module.MiddlewareHandler) {
	e := newExternalApplicationController()
	return apinto_module.RoutersInfo{{
			Method: http.MethodGet,
			Path:   "/api/external-apps",

			HandlerFunc: e.getList,
		},
			{
				Method: http.MethodGet,
				Path:   "/api/external-app",

				HandlerFunc: e.getInfo,
			},
			{
				Method: http.MethodPost,
				Path:   "/api/external-app",

				HandlerFunc: e.create,
			},
			{
				Method: http.MethodPut,
				Path:   "/api/external-app",

				HandlerFunc: e.edit,
			},
			{
				Method: http.MethodDelete,
				Path:   "/api/external-app",

				HandlerFunc: e.delete,
			},
			{
				Method: http.MethodPut,
				Path:   "/api/external-app/enable",

				HandlerFunc: e.enable,
			},
			{
				Method: http.MethodPut,
				Path:   "/api/external-app/disable",

				HandlerFunc: e.disable,
			},
			{
				Method: http.MethodPut,
				Path:   "/api/external-app/token",

				HandlerFunc: e.flushToken,
			},
		}, []apinto_module.MiddlewareHandler{
			middleware.CreateF(e.openApiCheck, func(api pm3.ApiInfo) bool {
				if api.Authority == pm3.Internal {
					return false
				}
				return strings.HasPrefix(api.Path, "/api2/")
			}),
		}
}
