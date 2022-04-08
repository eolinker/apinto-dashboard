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
	Name       string            `yaml:"name"`
	I18nNames  map[string]string `yaml:"i18n_name"`
	Titles     map[string]string `yaml:"titles"`
	Fields 	   []string `yaml:"fields"`
	Profession string            `yaml:"professions"`
}

type Config struct {
	Zone        string                  `yaml:"zone"`
	Default     string                  `yaml:"default"`
	Professions []*ProfessionConfigItem `yaml:"professions"`
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
	return setDefault(c), nil
}
func setDefault(c*Config) *Config  {
	if c.Zone == "" {
		c.Zone = string(apinto.ZhCn)
	}
	if c.Default == "" && len(c.Professions) > 0 {
		c.Default = c.Professions[0].Name
	}

	for _,p:=range c.Professions{
		if p.Titles[]
	}

}
func toModule(c *Config) []*apinto.Module {
	r := make([]*apinto.Module, 0, len(c.Professions))
	for _, cm := range c.Professions {
		m := &apinto.Module{
			Path:     fmt.Sprintf("/%s/list", cm.Name),
			Handler:  professions.NewProfession(cm.Name, cm.Profession, nil, nil),
			Name:     cm.Name,
			I18nName: make(map[apinto.ZoneName]string),
		}
		for k, v := range cm.I18nNames {
			m.I18nName[apinto.ZoneName(strings.ToLower(k))] = v
		}
		r = append(r, m)
	}
	return r
}
