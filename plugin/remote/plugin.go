package remote

import (
	"github.com/eolinker/apinto-dashboard/common"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

func checkConfig(config pm3.PluginConfig) (*Config, error) {

	if config == nil {
		return &Config{}, nil
	}
	c := &Config{}
	if ls, has := config["header"]; has {
		c.Header = common.SliceToMapO(ls, func(t pm3.ExtendParams) (string, string) {
			return t.Name, t.Value
		})
	}
	if ls, has := config["query"]; has {
		c.Query = common.SliceToMapO(ls, func(t pm3.ExtendParams) (string, string) {
			return t.Name, t.Value
		})
	}
	if ls, has := config["initialize"]; has {
		c.Initialize = common.SliceToMapO(ls, func(t pm3.ExtendParams) (string, string) {
			return t.Name, t.Value
		})
	}
	return c, nil
}

func DecodeDefine(define interface{}) (*Define, error) {
	return apinto_module.DecodeFor[Define](define)
}
