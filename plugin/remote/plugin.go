package remote

import (
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

type rPlugin struct {
	define *Define
}

func newRPlugin(define interface{}) (apinto_module.Plugin, error) {
	dv, err := DecodeDefine(define)
	if err != nil {
		return nil, err
	}
	return &rPlugin{define: dv}, nil
}

func (r *rPlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	cv, err := r.checkConfig(name, config)
	if err != nil {
		return nil, err
	}

	module := &rModule{name: name}

	//注册远程插件存储接口, 打开方式接口
	remoteStorage := newRemotePluginController(name, cv, r.define)
	module.routersInfo = append(module.routersInfo, remoteStorage.createRemoteApis()...)

	return module, nil
}

func (r *rPlugin) CheckConfig(name string, config interface{}) error {
	_, err := r.checkConfig(name, config)
	if err != nil {
		return err
	}
	return nil
}
func (r *rPlugin) checkConfig(name string, config interface{}) (*Config, error) {
	return DecodeConfig(config)
}

func (r *rPlugin) GetPluginFrontend(moduleName string) string {
	return fmt.Sprintf("remote/%s", moduleName)
}

func (r *rPlugin) IsPluginVisible() bool {
	return true
}

func (r *rPlugin) IsShowServer() bool {
	return !r.define.Internet
}

func (r *rPlugin) IsCanUninstall() bool {
	return true
}

func (r *rPlugin) IsCanDisable() bool {
	return true
}

func DecodeDefine(define interface{}) (*Define, error) {
	return apinto_module.DecodeFor[Define](define)
}
func DecodeConfig(config interface{}) (*Config, error) {
	return apinto_module.DecodeFor[Config](config)

}
