package main

import (
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type ProfessionConfigItem struct {
	Name       string              `yaml:"name"`
	I18nNames  map[string]string   `yaml:"i18n_name"`
	Titles     map[string][]string `yaml:"titles"`
	Fields     []string            `yaml:"fields"`
	Profession string              `yaml:"profession"`
}

type UserDetailsConfig struct {
	Type string `yaml:"type"`
	File string `yaml:"file"`
}

type Config struct {
	Zone            string                  `yaml:"zone"`
	Default         string                  `yaml:"default"`
	Apinto          []string                `yaml:"apinto"`
	Port            string                  `yaml:"port"`
	FilterForwarded bool                    `yaml:"filter_forwarded"`
	UserDetails     *UserDetailsConfig      `yaml:"user_details"`
	Professions     []*ProfessionConfigItem `yaml:"professions"`
}

func ReadConfig(file string) (*Config, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	if c.Zone == "" {
		c.Zone = string(apinto.ZhCn)
	}
	if c.Default == "" && len(c.Professions) > 0 {
		c.Default = c.Professions[0].Name
	}
	return c, nil
}

var (
	defaultTitles = map[apinto.ZoneName][]string{
		apinto.ZhCn: {"名称", "驱动", "创建时间", "更新时间"},
		apinto.EnUs: {"name", "driver", "create time", "update time"},
	}
	defaultFields = []string{"name", "driver", "create", "update"}
)

func titleView(titles map[string][]string, fields []string) (map[apinto.ZoneName][]string, []string) {

	if titles == nil && len(titles) == 0 {
		return defaultTitles, defaultFields
	}
	titlesMap := make(map[apinto.ZoneName][]string)
	for name, list := range titles {
		titles[strings.ToLower(name)] = list
	}
	for zn, dv := range defaultTitles {
		if v, has := titles[string(zn)]; !has || len(v) == 0 {
			titlesMap[zn] = dv
		} else {
			titlesMap[zn] = v
		}
	}
	if len(fields) == 0 {
		fields = defaultFields
	}
	return titlesMap, fields
}
func toModule(c *Config) []*apinto.Module {
	r := make([]*apinto.Module, 0, len(c.Professions))
	for _, cm := range c.Professions {
		titles, fields := titleView(cm.Titles, cm.Fields)
		m := &apinto.Module{
			Path:     fmt.Sprintf("/%s/list", cm.Name),
			Handler:  professions.NewProfession(cm.Name, cm.Profession, titles, fields, nil),
			Name:     cm.Name,
			I18nName: make(map[apinto.ZoneName]string),
			NotView:  false,
		}
		for k, v := range cm.I18nNames {
			m.I18nName[apinto.ZoneName(strings.ToLower(k))] = v
		}
		r = append(r, m)
		r = append(r, &apinto.Module{
			NotView: true,
			Handler: m.Handler,
			Path:    fmt.Sprintf("/profession/%s/", cm.Name),
		})
	}
	return r
}
