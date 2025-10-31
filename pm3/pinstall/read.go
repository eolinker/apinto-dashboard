package pinstall

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/pm3"
	"gopkg.in/yaml.v3"
)

type pluginInstallCompatible struct {
	*pm3.PluginDefine
	Front      string `json:"front,omitempty" yaml:"front"`
	Navigation string `json:"navigation,omitempty" yaml:"navigation"`
}

func Read(input []byte) (*pm3.PluginDefine, error) {

	p := new(pluginInstallCompatible)

	err := yaml.Unmarshal(input, &p.PluginDefine)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(input, p)
	if err != nil {
		return nil, err
	}

	return transform(p), nil
}

type pluginInstallInfoInner struct {
	Auto         bool `json:"auto,omitempty" yaml:"auto"`
	IsCanDisable bool `json:"is_can_disable,omitempty" yaml:"is_can_disable"`
}

func ReadInner(input []byte) (info *pm3.PluginDefine, auto bool, isCanDisable bool, err error) {
	define, err := Read(input)
	if err != nil {
		return nil, false, false, err
	}
	p := new(pluginInstallInfoInner)
	err = yaml.Unmarshal(input, p)
	if err != nil {
		return nil, false, false, err
	}

	return define, p.Auto, p.IsCanDisable, nil
}

func transform(info *pluginInstallCompatible) *pm3.PluginDefine {
	if info.Version == "" {
		info.Version = "v0.0.0"
	}

	if len(info.Navigations) == 0 && info.Navigation != "" {
		info.Navigations = []pm3.NavigationItem{
			{
				Navigation: info.Navigation,
				Router:     info.Front,
				Name:       fmt.Sprintf("%s", info.Name),
				Cname:      info.Cname,
				Access: []pm3.AccessItem{
					{
						Name:  fmt.Sprintf("%s.%s.view", info.Id, info.Name),
						Cname: "查看",
					},
					{
						Name:   fmt.Sprintf("%s.%s.edit", info.Id, info.Name),
						Cname:  "编辑",
						Depend: []string{fmt.Sprintf("%s.%s.view", info.Id, info.Name)},
					},
				},
			},
		}
	}

	return info.PluginDefine
}
