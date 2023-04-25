package plugin

import (
	_ "embed"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"gopkg.in/yaml.v3"
)

//go:embed apinto_plugin.yml
var pluginData []byte

var pluginConf []*v1.GlobalPlugin

func init() {
	var err error

	pc := make([]*v1.GlobalPlugin, 0)
	err = yaml.Unmarshal(pluginData, &pc)
	if err != nil {
		panic(err)
	}

	pluginConf = pc
}

func GetGlobalPluginConf() []*v1.GlobalPlugin {
	return pluginConf
}
