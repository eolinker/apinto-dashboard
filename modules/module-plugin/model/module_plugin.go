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
	UUID  string
	Name  string
	Count int
}

type ModulePluginInfo struct {
	*entry.ModulePlugin
	Enable     bool
	CanDisable bool
	Uninstall  bool
}

type PluginEnableInfo struct {
	Name       string
	Server     string
	Header     []*ExtendParams
	Query      []*ExtendParams
	Initialize []*ExtendParams
}

type PluginEnableRender struct {
	Internet     bool
	NameConflict bool
	//Invisible  bool
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
	Config interface{}
	Define interface{}
}

// PluginEnableCfg 插件启用配置
type PluginEnableCfg struct {
	Server     string          `json:"server"`
	Header     []*ExtendParams `json:"header"`
	Query      []*ExtendParams `json:"query"`
	Initialize []*ExtendParams `json:"initialize"`
}

// PluginEnableCfgMap 创建插件时启用的配置
type PluginEnableCfgMap struct {
	Server     string            `json:"server"`
	Header     map[string]string `json:"header"`
	Query      map[string]string `json:"query"`
	Initialize map[string]string `json:"initialize"`
}

// PluginDefine 插件安装文件里的Define配置
type PluginDefine struct {
	Internet   bool                 `json:"internet,omitempty" yaml:"internet,omitempty"`     //remote
	Server     string               `json:"server,omitempty" yaml:"server,omitempty"`         //remote
	Path       string               `json:"path,omitempty" yaml:"path,omitempty"`             //remote
	Middleware []*MiddlewareDefine  `json:"middleware,omitempty" yaml:"middleware,omitempty"` //local
	Router     *RouterDefine        `json:"router,omitempty" yaml:"router,omitempty"`         //local
	Headers    []ExtendParamsRender `json:"headers,omitempty" yaml:"headers,omitempty"`       //local
	Querys     []ExtendParamsRender `json:"querys,omitempty" yaml:"querys,omitempty"`         //remote local
	Initialize []ExtendParamsRender `json:"initialize,omitempty" yaml:"initialize,omitempty"` //remote local
	Profession string               `json:"profession" yaml:"profession"`                     //dynamic
	Skill      string               `json:"skill" yaml:"skill"`                               //dynamic
	Drivers    []DynamicTitleDefine `json:"drivers" yaml:"drivers"`                           //dynamic
	Fields     []DynamicTitleDefine `json:"fields" yaml:"fields"`                             //dynamic
	Render     map[string]string    `json:"render" yaml:"render"`                             //dynamic
}

type MiddlewareDefine struct {
	Name string   `json:"name" yaml:"name"`
	Path string   `json:"path" yaml:"path"`
	Life string   `json:"life" yaml:"life"`
	Rule []string `json:"rule" yaml:"rule"`
}

type RouterDefine struct {
	Home     string                               `json:"home" yaml:"home"`
	Html     []*RouterHtmlDefine                  `json:"html" yaml:"html"`
	Frontend []string                             `json:"frontend" yaml:"frontend"`
	Api      map[string]map[string]*ApiPathDefine `json:"api" yaml:"api"`
	Openapi  map[string]map[string]*ApiPathDefine `json:"openapi" yaml:"openapi"`
}

type RouterHtmlDefine struct {
	Path  string   `json:"path" yaml:"path"`
	Label []string `json:"label" yaml:"label"`
}

type ApiPathDefine struct {
	Label []string `json:"label" yaml:"label"`
}

type DynamicTitleDefine struct {
	Name  string `json:"name" yaml:"name"`
	Title string `json:"title" yaml:"title"`
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
	ID         string        `json:"id" yaml:"id"`
	Name       string        `json:"name" yaml:"name"`
	Version    string        `json:"version" yaml:"version"`
	CName      string        `json:"cname" yaml:"cname"`
	Navigation string        `json:"navigation" yaml:"navigation"`
	GroupID    string        `json:"group_id" yaml:"group_id"`
	Resume     string        `json:"resume" yaml:"resume"`
	ICon       string        `json:"icon" yaml:"icon"`
	Driver     string        `json:"driver" yaml:"driver"`
	Define     *PluginDefine `json:"define" yaml:"define"`
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
	GroupID    string `json:"group_id" yaml:"group_id"`
	Type       int    `json:"type" yaml:"type"`
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
	UUID  string
	CName string
}

// NavigationModuleInfo 导航所需要的模块信息
type NavigationModuleInfo struct {
	Name       string
	Title      string
	Type       string
	Path       string
	Navigation string
}
