package application_controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	c := newApplicationController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/applications",
			HandlerFunc: c.lists,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/enum",

			HandlerFunc: c.enum,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/application",

			HandlerFunc: c.createApp,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application",

			HandlerFunc: c.info,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/application",

			HandlerFunc: c.updateApp,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/application",

			HandlerFunc: c.deleteApp,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/onlines",

			HandlerFunc: c.getOnlineInfo,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/application/online",

			HandlerFunc: c.online,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/application/offline",

			HandlerFunc: c.offline,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/drivers",

			HandlerFunc: c.drivers,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/auths",

			HandlerFunc: c.auths,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/auth",

			HandlerFunc: c.getAuth,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/application/auth/details",

			HandlerFunc: c.getAuthDetails,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/application/auth",

			HandlerFunc: c.createAuth,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/application/auth",

			HandlerFunc: c.updateAuth,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/application/auth",

			HandlerFunc: c.delAuth,
		},
	}
}
