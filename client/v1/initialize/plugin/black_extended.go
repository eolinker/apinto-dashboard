package plugin

import (
	_ "embed"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/eosc/log"
	"gopkg.in/yaml.v3"
	"os"
)

const blackExtendedYamlPath = "./black_extended.yml"

//go:embed black_extended.yml
var blackExtendedData []byte

var blackExtendedConf []*v1.GlobalPlugin

func init() {
	var data []byte
	var err error

	data, err = os.ReadFile(blackExtendedYamlPath)
	if err != nil {
		log.Info("black_extended.yml doesn't exist. Read embedded data. ")
		//文件不存在则读取内嵌文件
		data = blackExtendedData
	}

	pc := make([]*v1.GlobalPlugin, 0)
	err = yaml.Unmarshal(data, &pc)
	if err != nil {
		panic(err)
	}

	blackExtendedConf = pc
}

func GetBlackExtendedPluginConf() []*v1.GlobalPlugin {
	return blackExtendedConf
}
