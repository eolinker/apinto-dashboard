package initialize

import (
	"embed"
	"fmt"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/embed_registry"
	"io/fs"
	"net/http"
	"path"

	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"gopkg.in/yaml.v3"

	"github.com/eolinker/eosc/log"
)

var (
	//go:embed plugins
	pluginDir embed.FS
)

func init() {
	plugins, err := loadPlugins("plugins", "plugin.yml")
	if err != nil {
		panic(err)
	}
	embed_registry.RegisterEmbedPlugin(plugins...)
}

func loadPlugins(dir string, target string) ([]*model.EmbedPluginCfg, error) {
	entries, err := pluginDir.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	plugins := make([]*model.EmbedPluginCfg, 0)
	for _, e := range entries {
		filePath := path.Join(dir, e.Name(), target)
		fileContent, err := pluginDir.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		pluginCfg := new(model.InnerPluginCfg)
		err = yaml.Unmarshal(fileContent, pluginCfg)
		if err != nil {
			log.Errorf("parse file(%s) error: %v", filePath, err)
			return nil, err
		}

		plugins = append(plugins, &model.EmbedPluginCfg{
			PluginCfg: pluginCfg,
			Resources: &model.EmbedPluginResources{
				PluginID: pluginCfg.ID,
				Icon:     pluginCfg.Icon,
				Readme:   "README.md",
				Fs:       pluginDir,
			},
		})

	}
	return plugins, nil
}

func GetInnerPluginFS(filePath string) (http.FileSystem, error) {
	//先检验文件是否存在
	_, err := pluginDir.ReadFile(fmt.Sprintf("plugins/%s", filePath))
	if err != nil {
		return nil, err
	}

	pluginsFS, err := fs.Sub(pluginDir, "plugins")
	if err != nil {
		return nil, err
	}
	return http.FS(pluginsFS), nil
}
