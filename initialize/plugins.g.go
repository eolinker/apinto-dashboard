package initialize

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"

	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"

	"gorm.io/gorm"

	"github.com/eolinker/eosc/log"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
)

var (
	//go:embed plugins
	pluginDir embed.FS
)

type Plugin struct {
	Id      string         `yaml:"id"`
	Name    string         `yaml:"name"`
	CName   string         `yaml:"cname"`
	Version string         `yaml:"version"`
	Icon    string         `yaml:"icon"`
	Driver  string         `yaml:"driver"`
	Core    bool           `yaml:"core"`
	Install *PluginInstall `yaml:"install"`
	Org     string         `yaml:"-"`
}

type PluginInstall struct {
	Auto       bool   `yaml:"auto"`
	Front      string `yaml:"front"`
	Navigation string `yaml:"navigation"`
}

func initPlugins(service module_plugin.IModulePluginService) error {
	ctx := context.Background()
	plugins, err := loadPlugins("plugins", "plugin.yml")
	if err != nil {
		return err
	}
	for _, p := range plugins {
		_, err := service.GetPluginInfo(ctx, p.Id)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			// 插入安装记录
			err = service.InstallPlugin(ctx, 0, "", &model.PluginYmlCfg{
				ID:     p.Id,
				Name:   p.Name,
				CName:  p.CName,
				Resume: "",
				ICon:   p.Icon,
				Driver: p.Driver,
			}, []byte(""))
			if err != nil {
				return err
			}
			navigation := ""
			if p.Install != nil {
				navigation = p.Install.Navigation
			}
			err = service.EnablePlugin(ctx, 0, p.Id, &dto.PluginEnableInfo{
				Name:       p.Name,
				Navigation: navigation,
				ApiGroup:   "",
			})
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func loadPlugins(dir string, target string) ([]*Plugin, error) {
	entry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	plugins := make([]*Plugin, 0)
	for _, e := range entry {
		nextFile := fmt.Sprintf("%s/%s", dir, e.Name())
		if e.IsDir() {
			s, err := loadPlugins(nextFile, target)
			if err != nil {
				return nil, err
			}
			plugins = append(plugins, s...)
			continue
		}
		if e.Name() == target {
			s, err := os.ReadFile(nextFile)
			if err != nil {
				return nil, err
			}
			var p Plugin
			err = json.Unmarshal(s, &p)
			if err != nil {
				log.Errorf("parse file(%s) error: %w")
			}
			p.Org = string(s)
			plugins = append(plugins, &p)
		}
	}
	return plugins, nil
}
