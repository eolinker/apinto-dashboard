package remote

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

type rModule struct {
	name              string
	middlewareHandler []apinto_module.MiddlewareHandler
	routersInfo       apinto_module.RoutersInfo
}

func (m *rModule) RoutersInfo() apinto_module.RoutersInfo {
	return m.routersInfo
}

func (m *rModule) Name() string {
	return m.name
}

func (m *rModule) Routers() (apinto_module.Routers, bool) {
	if len(m.routersInfo) == 0 {
		return nil, false
	}
	return m, true
}

func (m *rModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func (m *rModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}
