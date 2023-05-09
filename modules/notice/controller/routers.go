package controller

import (
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func InitRouter() apinto_module.RoutersInfo {
	c := newNoticeController()
	return []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/channels",
			Handler:     "notice.getChannelsEnum",
			HandlerFunc: []apinto_module.HandlerFunc{c.channels},
		},
	}
}
