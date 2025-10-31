package local

import (
	"github.com/eolinker/apinto-dashboard/common"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

func (d *tDriver) createModule(id, name string, define *Define, config pm3.PluginConfig) (apinto_module.Module, error) {
	cv, err := d.checkConfig(config)
	if err != nil {
		return nil, err
	}

	proxy, err := NewProxyAPi(define.Cmd, id, name, cv)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}

func (d *tDriver) checkConfig(config pm3.PluginConfig) (c *Config, err error) {
	c = &Config{}
	if config == nil {
		return
	}
	initializeParameters, has := config["initialize"]
	if !has {
		return
	}
	c.Initialize = common.SliceToMapO(initializeParameters, func(t pm3.ExtendParams) (string, string) {
		return t.Name, t.Value
	})

	return
}

func DecodeDefine(define interface{}) (*Define, error) {
	return apinto_module.DecodeFor[Define](define)
}
func DecodeConfig(config interface{}) (*Config, error) {
	return apinto_module.DecodeFor[Config](config)

}
