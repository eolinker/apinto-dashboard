package main

import (
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type Plugin struct {
	frontends   []pm3.FrontendAsset
	apis        []pm3.Api
	middlewares []shared.Middleware
}

func NewPlugin() *Plugin {
	service := Service{}

	apis := make([]pm3.Api, 0)
	apis = append(apis, pm3.Api{

		Authority: pm3.Public,
		Access:    "test.test.test",
		Method:    http.MethodPost,
		Path:      "/api/myservice/{id}",

		HandlerFunc: service.Test,
	})

	return &Plugin{
		apis: apis,
	}
}

func (p *Plugin) Frontend() []pm3.FrontendAsset {
	return []pm3.FrontendAsset{}
}

func (p *Plugin) Apis() []pm3.Api {
	return p.apis
}

func (p *Plugin) Middleware() []shared.Middleware {
	//TODO implement me
	panic("implement me")
}
