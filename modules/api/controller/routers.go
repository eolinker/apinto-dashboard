package controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	group_controller "github.com/eolinker/apinto-dashboard/modules/group/group-controller"
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	c := newApiController()
	g := group_controller.NewCommonGroupController()
	return []apinto_module.RouterInfo{

		{
			Method:      http.MethodGet,
			Path:        "/api/routers",
			Handler:     "api.routers",
			HandlerFunc: []apinto_module.HandlerFunc{c.routers},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/router",
			Handler:     "api.getInfo",
			HandlerFunc: []apinto_module.HandlerFunc{c.getInfo},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/router",
			Handler:     "api.create",
			HandlerFunc: []apinto_module.HandlerFunc{c.create},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/router",
			Handler:     "api.update",
			HandlerFunc: []apinto_module.HandlerFunc{c.update},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/router",
			Handler:     "api.delete",
			HandlerFunc: []apinto_module.HandlerFunc{c.delete},
		},

		{
			Method:      http.MethodPost,
			Path:        "/api/routers/batch-online",
			Handler:     "api.batchOnline",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.batchOnline},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/routers/batch-offline",
			Handler:     "api.batchOffline",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.batchOffline},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/routers/batch-online/check",
			Handler:     "api.batchOnlineCheck",
			HandlerFunc: []apinto_module.HandlerFunc{c.batchOnlineCheck},
		},

		{
			Method:      http.MethodPut,
			Path:        "/api/router/online",
			Handler:     "api.online",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.online},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/router/offline",
			Handler:     "api.offline",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.offline},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/router/online/info",
			Handler:     "api.getOnlineInfo",
			HandlerFunc: []apinto_module.HandlerFunc{c.getOnlineInfo},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/router/groups",
			Handler:     "api.groups",
			HandlerFunc: []apinto_module.HandlerFunc{c.groups},
		},

		{
			Method:      http.MethodGet,
			Path:        "/api/router/source",
			Handler:     "api.getSourceList",
			HandlerFunc: []apinto_module.HandlerFunc{c.getSourceList},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/router/import",
			Handler:     "api.getImportCheckList",
			HandlerFunc: []apinto_module.HandlerFunc{c.getImportCheckList},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/router/enum",
			Handler:     "api.routerEnum",
			HandlerFunc: []apinto_module.HandlerFunc{c.routerEnum},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/router/import",
			Handler:     "api.importAPI",
			HandlerFunc: []apinto_module.HandlerFunc{c.importAPI},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/group/:group_type",
			Handler:     "group.groups",
			HandlerFunc: []apinto_module.HandlerFunc{c.groups},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/group/:group_type",
			Handler:     "group.createGroup",
			HandlerFunc: []apinto_module.HandlerFunc{g.CreateGroup},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/group/:group_type/:uuid",
			Handler:     "group.updateGroup",
			HandlerFunc: []apinto_module.HandlerFunc{g.UpdateGroup},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/group/:group_type/:uuid",
			Handler:     "group.delGroup",
			HandlerFunc: []apinto_module.HandlerFunc{g.DelGroup},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/groups/:group_type/sort",
			Handler:     "group.groupSort",
			HandlerFunc: []apinto_module.HandlerFunc{g.GroupSort},
		},
	}
}
