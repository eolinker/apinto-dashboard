package remote

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type rModule struct {
	*pm3.ModuleTool

	name              string
	middlewareHandler []apinto_module.MiddlewareHandler
	routersInfo       apinto_module.RoutersInfo
}

func (m *rModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (m *rModule) Apis() []pm3.Api {
	return m.routersInfo
}

func (m *rModule) Middleware() []pm3.Middleware {
	return m.middlewareHandler
}

func (m *rModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (m *rModule) Name() string {
	return m.name
}
