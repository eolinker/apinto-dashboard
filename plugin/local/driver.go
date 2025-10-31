package local

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type tDriver struct {
}

func (d *tDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *tDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	define, err := DecodeDefine(info.Define)
	if err != nil {
		return nil, err
	}
	return d.createModule(info.Id, info.Name, define, config)

}

func NewDriver() apinto_module.Driver {
	return &tDriver{}
}
