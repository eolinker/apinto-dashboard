package controller

import apinto_module "github.com/eolinker/apinto-module"

var (
	_ apinto_module.Driver = (*Plugin)(nil)
	_ apinto_module.Plugin = (*Plugin)(nil)
	_ apinto_module.Module = (*Module)(nil)
)

type Plugin struct {
}

func (p *Plugin) CreateModule(name string, apiPrefix string, config interface{}) (apinto_module.Module, error) {
	return NewModule(), nil
}

func (p *Plugin) CheckConfig(name string, apiPrefix string, config interface{}) error {
	return nil
}

func (p *Plugin) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return p, nil
}

type Module struct {
	name string
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Routers() (apinto_module.Routers, bool) {
	return nil, false
}

func (m *Module) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func (m *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func NewModule() *Module {
	return &Module{}
}
