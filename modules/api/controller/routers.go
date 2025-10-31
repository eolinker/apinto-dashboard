package controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	group_controller "github.com/eolinker/apinto-dashboard/modules/group/group-controller"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

func initRouter(name string) apinto_module.RoutersInfo {
	c := newApiController()
	g := group_controller.NewCommonGroupController()
	return []apinto_module.RouterInfo{

		{
			Method: http.MethodGet,
			Path:   "/api/routers",

			HandlerFunc: c.routers,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/router",

			HandlerFunc: c.getInfo,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/router",

			HandlerFunc: c.create,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/router",

			HandlerFunc: c.update,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/router",

			HandlerFunc: c.delete,
		},

		{
			Method: http.MethodPost,
			Path:   "/api/routers/batch-online",

			HandlerFunc: c.batchOnline,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/routers/batch-offline",

			HandlerFunc: c.batchOffline,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/routers/batch-online/check",

			HandlerFunc: c.batchOnlineCheck,
		},

		{
			Method: http.MethodPut,
			Path:   "/api/router/online",

			HandlerFunc: c.online,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/router/offline",

			HandlerFunc: c.offline,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/router/online/info",

			HandlerFunc: c.getOnlineInfo,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/router/groups",
			Authority:   pm3.Public,
			HandlerFunc: c.groups,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/router/source",

			HandlerFunc: c.getSourceList,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/router/import",

			HandlerFunc: c.getImportCheckList,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/router/enum",

			HandlerFunc: c.routerEnum,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/router/import",

			HandlerFunc: c.importAPI,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/router/check",

			HandlerFunc: c.checkApiExist,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/group/:group_type",

			HandlerFunc: c.groups,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/group/:group_type",

			HandlerFunc: g.CreateGroup,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/group/:group_type/:uuid",

			HandlerFunc: g.UpdateGroup,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/group/:group_type/:uuid",

			HandlerFunc: g.DelGroup,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/groups/:group_type/sort",

			HandlerFunc: g.GroupSort,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/group/:group_type/check",

			HandlerFunc: g.CheckGroupExist,
		},
	}
}
