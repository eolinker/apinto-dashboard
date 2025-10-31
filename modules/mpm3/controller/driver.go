package controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type PM3Module struct {
	*pm3.ModuleTool

	apis []pm3.Api
}

var (
	_ apinto_module.Driver = (*Driver)(nil)
)

type Driver struct {
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *Driver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewPM3Module(info.Id, info.Name), nil
}

func NewPM3Module(id, name string) *PM3Module {
	apis := make([]pm3.Api, 0, 20)
	apis = append(apis, NewInstall().apis()...)
	apis = append(apis, NewFrontendController().apis()...)
	apis = append(apis, NewPlugin().Apis()...)
	m := &PM3Module{
		ModuleTool: pm3.NewModuleTool(id, name),
	}
	m.InitAccess(apis)
	m.apis = apis
	return m
}

func (P *PM3Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (P *PM3Module) Apis() []pm3.Api {

	return P.apis
}

func (P *PM3Module) Middleware() []pm3.Middleware {
	return nil
}

func (P *PM3Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}
