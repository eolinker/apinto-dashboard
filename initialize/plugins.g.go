package initialize

import (
	"context"
	"embed"
	"fmt"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc/common/bean"
	"gopkg.in/yaml.v3"
	"net/http"
	"path"

	"gorm.io/gorm"

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
	Core       bool   `yaml:"core"`
	Auto       bool   `yaml:"auto"`
	Front      string `yaml:"front"`
	Navigation string `yaml:"navigation"`
}

func InitPlugins() error {
	var service module_plugin.IModulePluginService
	bean.Autowired(&service)
	ctx := context.Background()
	plugins, err := loadPlugins("plugins", "plugin.yml")
	if err != nil {
		return err
	}
	for _, p := range plugins {
		//TODO 校验内置插件

		pluginInfo, err := service.GetPluginInfo(ctx, p.Id)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// 插入安装记录
			err = service.InstallInnerPlugin(ctx, &model.InnerPluginYmlCfg{
				ID:         p.Id,
				Name:       p.Name,
				Version:    p.Version,
				CName:      p.CName,
				Resume:     p.Resume,
				ICon:       p.Icon,
				Driver:     p.Driver,
				Front:      p.Front,
				Navigation: p.Navigation,
				Core:       p.Core,
				Auto:       p.Auto,
			})
			if err != nil {
				return err
			}
			continue
		}
		//TODO 判断version有没改变，有则更新
		if pluginInfo.Version != p.Version {

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

func GetInnerPluginFSHandler(stripPrefix, filePath string) (http.Handler, error) {
	fileServer := http.StripPrefix(stripPrefix, http.FileServer(http.FS(pluginDir)))
	_, err := pluginDir.ReadFile(fmt.Sprintf("plugins/%s", filePath))
	return fileServer, err
}
