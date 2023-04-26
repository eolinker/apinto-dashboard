package initialize

import (
	"context"
	"embed"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"io/fs"
	"net/http"
	"path"

	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc/common/bean"
	"gopkg.in/yaml.v3"

	"github.com/eolinker/eosc/log"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
)

var (
	//go:embed plugins
	pluginDir embed.FS
)

type Plugin struct {
	Id         string `yaml:"id"`
	Name       string `yaml:"name"`
	CName      string `yaml:"cname"`
	Resume     string `yaml:"resume"`
	Version    string `yaml:"version"`
	Icon       string `yaml:"icon"`
	Driver     string `yaml:"driver"`
	Front      string `yaml:"front"`
	Navigation string `yaml:"navigation"`
	GroupID    string `yaml:"group_id"`
	Type       int    `yaml:"type"`
	Auto       bool   `yaml:"auto"`
}

func InitPlugins() error {
	var service module_plugin.IModulePluginService
	bean.Autowired(&service)
	ctx := context.Background()

	plugins, err := loadPlugins("plugins", "plugin.yml")
	if err != nil {
		return err
	}

	innerPlugins, err := service.GetInnerPluginList(ctx)
	if err != nil {
		return err
	}
	innerPluginsMap := common.SliceToMap(innerPlugins, func(t *model.ModulePluginInfo) string {
		return t.UUID
	})
	for _, p := range plugins {
		//TODO 校验内置插件

		pluginCfg := &model.InnerPluginYmlCfg{
			ID:         p.Id,
			Name:       p.Name,
			Version:    p.Version,
			CName:      p.CName,
			Resume:     p.Resume,
			ICon:       p.Icon,
			Driver:     p.Driver,
			Front:      p.Front,
			Navigation: p.Navigation,
			GroupID:    p.GroupID,
			Type:       p.Type,
			Auto:       p.Auto,
		}

		pluginInfo, has := innerPluginsMap[p.Id]
		if !has {
			// 插入安装记录
			err = service.InstallInnerPlugin(ctx, pluginCfg)
			if err != nil {
				return err
			}
			continue
		} else {
			//判断version有没改变，有则更新
			if pluginInfo.Version != p.Version {
				err = service.UpdateInnerPlugin(ctx, pluginCfg)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func loadPlugins(dir string, target string) ([]*Plugin, error) {
	entries, err := pluginDir.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	plugins := make([]*Plugin, 0)
	for _, e := range entries {
		nextFile := path.Join(dir, e.Name())
		if e.IsDir() {
			s, err := loadPlugins(nextFile, target)
			if err != nil {
				return nil, err
			}
			plugins = append(plugins, s...)
			continue
		}
		if e.Name() == target {
			s, err := pluginDir.ReadFile(nextFile)
			if err != nil {
				return nil, err
			}
			p := new(Plugin)
			err = yaml.Unmarshal(s, p)
			if err != nil {
				log.Errorf("parse file(%s) error: %w")
				return nil, err
			}
			plugins = append(plugins, p)
		}
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
