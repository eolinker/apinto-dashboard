package discovery_controller

import (
	apinto_module "github.com/eolinker/apinto-module"
)

func initRouter(name string) apinto_module.RoutersInfo {
	return apinto_module.RoutersInfo{}
	//c := newDiscoveryController()
	//return apinto_module.RoutersInfo{
	//	{
	//		Method:      http.MethodGet,
	//		Path:        "/api/discoveries",
	//		Handler:     "discovery.getList",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.getList},
	//	},
	//
	//	{
	//		Method:      http.MethodGet,
	//		Path:        "/api/discovery",
	//		Handler:     "discovery.getInfo",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.getInfo},
	//	},
	//
	//	{
	//		Method:      http.MethodPost,
	//		Path:        "/api/discovery",
	//		Handler:     "discovery.create",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.create},
	//	},
	//
	//	{
	//		Method:      http.MethodPut,
	//		Path:        "/api/discovery",
	//		Handler:     "discovery.update",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.update},
	//	},
	//
	//	{
	//		Method:      http.MethodDelete,
	//		Path:        "/api/discovery",
	//		Handler:     "discovery.delete",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.delete},
	//	},
	//
	//	{
	//		Method:      http.MethodGet,
	//		Path:        "/api/discovery/enum",
	//		Handler:     "discovery.getEnum",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.getEnum},
	//	},
	//
	//	{
	//		Method:      http.MethodGet,
	//		Path:        "/api/discovery/drivers",
	//		Handler:     "discovery.getDrivers",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.getDrivers},
	//	},
	//
	//	{
	//		Method:      http.MethodPut,
	//		Path:        "/api/discovery/:discovery_name/online",
	//		Handler:     "discovery.online",
	//		HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.online},
	//	},
	//
	//	{
	//		Method:      http.MethodPut,
	//		Path:        "/api/discovery/:discovery_name/offline",
	//		Handler:     "discovery.offline",
	//		HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, c.offline},
	//	},
	//
	//	{
	//		Method:      http.MethodGet,
	//		Path:        "/api/discovery/:discovery_name/onlines",
	//		Handler:     "discovery.getOnlineList",
	//		HandlerFunc: []apinto_module.HandlerFunc{c.getOnlineList},
	//	},
	//}
}
