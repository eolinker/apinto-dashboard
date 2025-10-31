package remote

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type rDriver struct {
}

func (d *rDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *rDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	define, err := DecodeDefine(info.Define)
	if err != nil {
		return nil, err
	}
	cv := new(Config)
	if len(config) > 0 {
		cv, err = checkConfig(config)
		if err != nil {
			return nil, err
		}
	}

	module := &rModule{name: info.Id, ModuleTool: pm3.NewModuleTool(info.Id, info.Name)}

	//注册远程插件存储接口, 打开方式接口
	remoteStorage := newRemotePluginController(info.Name, cv, define)
	module.routersInfo = append(module.routersInfo, remoteStorage.createRemoteApis()...)
	module.InitAccess(module.routersInfo)
	return module, nil
}

func NewDriver() apinto_module.Driver {
	return &rDriver{}
}
