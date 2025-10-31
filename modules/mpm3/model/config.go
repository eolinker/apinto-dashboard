package model

import "github.com/eolinker/apinto-dashboard/pm3"

type PluginConfig struct {
	UUID   string
	Driver string
	Define *pm3.PluginDefine
	Config []byte
}

// PluginEnableCfg 插件启用配置
type PluginEnableCfg struct {
	//Server     string          `json:"server"`
	Header     []*ExtendParams `json:"header"`
	Query      []*ExtendParams `json:"query"`
	Initialize []*ExtendParams `json:"initialize"`
}

type ExtendParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewExtendParams(name string, value string) ExtendParams {
	return ExtendParams{Name: name, Value: value}
}

type PluginEnableRender struct {
	//Internet     bool
	//Server string //define里的server
	//NameConflict bool
	//Invisible  bool
	Headers    []ExtendParamsRender
	Querys     []ExtendParamsRender
	Initialize []ExtendParamsRender
}
type ExtendParamsRender struct {
	Name        string `json:"name" yaml:"name"`
	Value       string `json:"value" yaml:"value"`
	Title       string `json:"title" yaml:"title"`
	Type        string `json:"type" yaml:"type"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Desc        string `json:"desc" yaml:"desc"`
}
