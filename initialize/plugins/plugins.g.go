/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package plugins

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
	//go:embed embed
	pluginDir embed.FS
)

func init() {
	plugins, err := loadPlugins("embed", "plugin.yml")
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
