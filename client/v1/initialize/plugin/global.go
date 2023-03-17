package plugin

import (
	_ "embed"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/eosc/log"
	"gopkg.in/yaml.v3"
	"os"
)

const pluginYamlPath = "./apinto_plugin.yml"

//go:embed apinto_plugin.yml
var pluginData []byte

var pluginConf []*v1.GlobalPlugin

func init() {
	var data []byte
	var err error

	data, err = os.ReadFile(pluginYamlPath)
	if err != nil {
		log.Info("apinto_plugin.yml doesn't exist. Read embedded data. ")
		//文件不存在则读取内嵌文件
		data = pluginData
	}

	pc := make([]*v1.GlobalPlugin, 0)
	err = yaml.Unmarshal(data, &pc)
	if err != nil {
		panic(err)
	}

	pluginConf = pc
}

func GetGlobalPluginConf() []*v1.GlobalPlugin {
	return pluginConf
}
