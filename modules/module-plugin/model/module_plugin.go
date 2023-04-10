package model

import (
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
)

type ModulePluginItem struct {
	*entry.ModulePlugin
	IsEnable bool
	IsInner  bool
}

type PluginGroup struct {
	UUID string
	Name string
}

type ModulePluginInfo struct {
	*entry.ModulePlugin
	IsEnable  bool
	Uninstall bool
}

type PluginEnableInfo struct {
	Name       string
	Navigation string
	ApiGroup   string
	Server     string
	Header     []*ExtendParams
	Query      []*ExtendParams
	Initialize []*ExtendParams
}

type PluginEnableRender struct {
	Internet   bool
	Invisible  bool
	ApiGroup   bool
	Headers    []ExtendParamsRender
	Querys     []ExtendParamsRender
	Initialize []ExtendParamsRender
}

type ExtendParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExtendParamsRender struct {
	Name        string `json:"name" yaml:"name"`
	Value       string `json:"value" yaml:"value"`
	Title       string `json:"title" yaml:"title"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Desc        string `json:"desc" yaml:"desc"`
}

type EnabledPlugin struct {
	UUID   string
	Name   string
	Driver string
	Config *PluginEnableCfg
	Define interface{}
}

// PluginEnableCfg 插件启用时的配置
type PluginEnableCfg struct {
	APIGroup   string          `json:"api_group"`
	Server     string          `json:"server"`
	Header     []*ExtendParams `json:"header"`
	Query      []*ExtendParams `json:"query"`
	Initialize []*ExtendParams `json:"initialize"`
}

// RemoteDefine 插件配置文件的driver为remote时的详细配置
type RemoteDefine struct {
	Internet   bool                 `json:"internet" yaml:"internet"`
	Server     string               `json:"server" yaml:"server"`
	Path       string               `json:"path" yaml:"path"`
	Querys     []ExtendParamsRender `json:"querys" yaml:"querys"`
	Initialize []ExtendParamsRender `json:"initialize" yaml:"initialize"`
}

// LocalDefine 插件配置文件的driver为local时的详细配置
type LocalDefine struct {
	Middleware []*MiddlewareDefine  `json:"middleware" yaml:"middleware"`
	ApiGroup   string               `json:"apigroup" yaml:"apigroup"`
	Api        *ApiDefine           `json:"api" yaml:"api"`
	Path       string               `json:"path" yaml:"path"`
	Invisible  bool                 `json:"invisible" yaml:"invisible"`
	Headers    []ExtendParamsRender `json:"headers" yaml:"headers"`
	Querys     []ExtendParamsRender `json:"querys" yaml:"querys"`
	Initialize []ExtendParamsRender `json:"initialize" yaml:"initialize"`
}

type MiddlewareDefine struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Desc string `json:"desc" yaml:"desc"`
}

type ApiDefine struct {
	Prefix string                               `json:"prefix" yaml:"prefix"`
	Paths  map[string]map[string]*ApiPathDefine `json:"paths" yaml:"paths"`
}

type ApiPathDefine struct {
	Access []string `json:"access" yaml:"access"`
}

// ProfessionDefine 插件配置文件的driver为profession时的详细配置
type ProfessionDefine struct {
	Type       string            `json:"type" yaml:"type"`
	Profession string            `json:"profession" yaml:"profession"`
	Drivers    []string          `json:"drivers" yaml:"drivers"`
	Driver     string            `json:"driver" yaml:"driver"`
	Name       string            `json:"name" yaml:"name"`
	Render     map[string]string `json:"render" yaml:"render"`
}

type InnerDefine struct {
	ID      string                 `json:"id" yaml:"id"`
	Name    string                 `json:"name" yaml:"name"`
	Version string                 `json:"version" yaml:"version"`
	ICon    string                 `json:"icon" yaml:"icon"`
	Driver  string                 `json:"driver" yaml:"driver"`
	Core    bool                   `json:"core" yaml:"core"`
	Install *InnerPluginYmlInstall `json:"install" yaml:"install"`
	Main    *InnerPluginYmlMain    `json:"main" yaml:"main"`
}

// MiddlewareItem 拦截器项结构体
type MiddlewareItem struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type PluginYmlCfg struct {
	ID         string            `json:"id" yaml:"id"`
	Name       string            `json:"name" yaml:"name"`
	CName      string            `json:"cname" yaml:"cname"`
	Resume     string            `json:"resume" yaml:"resume"`
	ICon       string            `json:"icon" yaml:"icon"`
	Driver     string            `json:"driver" yaml:"driver"`
	Remote     *RemoteDefine     `json:"remote" yaml:"remote"`
	Local      *LocalDefine      `json:"local" yaml:"local"`
	Profession *ProfessionDefine `json:"profession" yaml:"profession"`
}

type InnerPluginYmlCfg struct {
	ID         string `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	Version    string `json:"version" yaml:"version"`
	CName      string `json:"cname" yaml:"cname"`
	Resume     string `json:"resume" yaml:"resume"`
	ICon       string `json:"icon" yaml:"icon"`
	Driver     string `json:"driver" yaml:"driver"`
	Front      string `json:"front" yaml:"front"`
	Navigation string `json:"navigation" yaml:"navigation"`
	Core       bool   `json:"core" yaml:"core"`
	Auto       bool   `json:"auto" yaml:"auto"`
}

type InnerPluginYmlInstall struct {
	Auto       bool   `json:"auto" yaml:"auto"`
	Front      string `json:"front" yaml:"front"`
	Navigation string `json:"navigation" yaml:"navigation"`
}

type InnerPluginYmlMain struct {
	Middleware []string `json:"middleware" yaml:"middleware"`
}

// NavigationEnabledPlugin 用于给导航返回的
type NavigationEnabledPlugin struct {
	*entry.ModulePluginEnable
	UUID string
}
