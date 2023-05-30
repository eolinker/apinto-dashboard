package local

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

type tModule struct {
	name              string
	middlewareHandler []apinto_module.MiddlewareHandler
	routersInfo       apinto_module.RoutersInfo
}

func (m *tModule) RoutersInfo() apinto_module.RoutersInfo {
	return m.routersInfo
}

func (m *tModule) MiddlewaresInfo() []apinto_module.MiddlewareHandler {
	return m.middlewareHandler
}

func (m *tModule) Name() string {
	return m.name
}

func (m *tModule) Routers() (apinto_module.Routers, bool) {
	if len(m.routersInfo) == 0 {
		return nil, false
	}
	return m, true
}

func (m *tModule) Middleware() (apinto_module.Middleware, bool) {
	if len(m.middlewareHandler) == 0 {
		return nil, false
	}
	return m, true
}

func (m *tModule) Support() (apinto_module.ProviderSupport, bool) {

	return nil, false

}
