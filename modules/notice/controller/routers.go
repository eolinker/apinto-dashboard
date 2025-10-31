package controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

func InitRouter() apinto_module.RoutersInfo {
	c := newNoticeController()
	return []apinto_module.RouterInfo{
		{
			Method:    http.MethodGet,
			Path:      "/api/channels",
			Authority: pm3.Public,

			HandlerFunc: c.channels,
		},
	}
}
