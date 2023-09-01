package application_controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	c := newApplicationController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/applications",
			Handler:     "applications.lists",
			HandlerFunc: []apinto_module.HandlerFunc{c.lists},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/enum",
			Handler:     "applications.enum",
			HandlerFunc: []apinto_module.HandlerFunc{c.enum},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/application",
			Handler:     "applications.createApp",
			HandlerFunc: []apinto_module.HandlerFunc{c.createApp},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application",
			Handler:     "applications.info",
			HandlerFunc: []apinto_module.HandlerFunc{c.info},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/application",
			Handler:     "applications.updateApp",
			HandlerFunc: []apinto_module.HandlerFunc{c.updateApp},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/application",
			Handler:     "applications.deleteApp",
			HandlerFunc: []apinto_module.HandlerFunc{c.deleteApp},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/onlines",
			Handler:     "applications.onlines",
			HandlerFunc: []apinto_module.HandlerFunc{c.getOnlineInfo},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/application/online",
			Handler:     "applications.online",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.online},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/application/offline",
			Handler:     "applications.offline",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.offline},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/drivers",
			Handler:     "applications.drivers",
			HandlerFunc: []apinto_module.HandlerFunc{c.drivers},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/auths",
			Handler:     "applications.auths",
			HandlerFunc: []apinto_module.HandlerFunc{c.auths},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/auth",
			Handler:     "applications.getAuth",
			HandlerFunc: []apinto_module.HandlerFunc{c.getAuth},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/application/auth/details",
			Handler:     "applications.getAuthDetails",
			HandlerFunc: []apinto_module.HandlerFunc{c.getAuthDetails},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/application/auth",
			Handler:     "applications.createAuth",
			HandlerFunc: []apinto_module.HandlerFunc{c.createAuth},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/application/auth",
			Handler:     "applications.updateAuth",
			HandlerFunc: []apinto_module.HandlerFunc{c.updateAuth},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/application/auth",
			Handler:     "applications.delAuth",
			HandlerFunc: []apinto_module.HandlerFunc{c.delAuth},
		},
	}
}
