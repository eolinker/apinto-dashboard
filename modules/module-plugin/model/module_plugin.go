package model

import (
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
)

type ModulePluginItem struct {
	*entry.PluginListItem
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
	Server       string //define里的server
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
	Type        string `json:"type" yaml:"type"`
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
	//Internet   bool                 `json:"internet,omitempty" yaml:"internet,omitempty"`     //remote
	Server string `json:"server,omitempty" yaml:"server,omitempty"` //remote
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`     //remote
	//Middleware []*MiddlewareDefine  `json:"middleware,omitempty" yaml:"middleware,omitempty"` //local
	//Router     *RouterDefine        `json:"router,omitempty" yaml:"router,omitempty"`         //local
	Headers    []ExtendParamsRender `json:"headers,omitempty" yaml:"headers,omitempty"`       //local
	Querys     []ExtendParamsRender `json:"querys,omitempty" yaml:"querys,omitempty"`         //remote local
	Initialize []ExtendParamsRender `json:"initialize,omitempty" yaml:"initialize,omitempty"` //remote local
	//Profession string               `json:"profession" yaml:"profession"`                     //dynamic
	//Skill      string               `json:"skill" yaml:"skill"`                               //dynamic
	//Drivers    []DynamicTitleDefine `json:"drivers" yaml:"drivers"`                           //dynamic
	//Fields     []DynamicTitleDefine `json:"fields" yaml:"fields"`                             //dynamic
	//Render     map[string]string    `json:"render" yaml:"render"`                             //dynamic
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

// MiddlewareItem 拦截器项结构体
type MiddlewareItem struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type ExternPluginCfg struct {
	ID         string      `json:"id" yaml:"id"`
	Name       string      `json:"name" yaml:"name"`
	Version    string      `json:"version" yaml:"version"`
	CName      string      `json:"cname" yaml:"cname"`
	Navigation string      `json:"navigation" yaml:"navigation"`
	GroupID    string      `json:"group_id" yaml:"group_id"`
	Resume     string      `json:"resume" yaml:"resume"`
	ICon       string      `json:"icon" yaml:"icon"`
	Driver     string      `json:"driver" yaml:"driver"`
	Define     interface{} `json:"define" yaml:"define"`
}

type InnerPluginCfg struct {
	ID                  string      `yaml:"id" json:"id"`
	Name                string      `yaml:"name" json:"name"`
	CName               string      `yaml:"cname" json:"cname"`
	Resume              string      `yaml:"resume" json:"resume"`
	Version             string      `yaml:"version" json:"version"`
	Icon                string      `yaml:"icon" json:"icon"`
	Driver              string      `yaml:"driver" json:"driver"`
	Front               string      `yaml:"front" json:"front"`
	Navigation          string      `yaml:"navigation" json:"navigation"`
	GroupID             string      `yaml:"group_id" json:"group_id"`
	Type                int         `yaml:"type" json:"type"`
	Auto                bool        `yaml:"auto" json:"auto"`
	IsCanDisable        bool        `yaml:"is_can_disable" json:"is_can_disable"`
	IsCanUninstall      bool        `yaml:"is_can_uninstall" json:"is_can_uninstall"`
	VisibleInNavigation bool        `yaml:"visible_in_navigation" json:"visible_in_navigation"`
	VisibleInMarket     bool        `yaml:"visible_in_market" json:"visible_in_market"`
	Define              interface{} `yaml:"define" json:"define"`
}

type PluginCfg struct {
	Version    string      `json:"version" yaml:"version"`
	Navigation string      `json:"navigation" yaml:"navigation"`
	GroupID    string      `json:"group_id" yaml:"group_id"`
	Resume     string      `json:"resume" yaml:"resume"`
	Type       int         `json:"type" yaml:"type"`
	Define     interface{} `json:"define" yaml:"define"`
}

type EmbedPluginCfg struct {
	PluginCfg *InnerPluginCfg
	Resources *EmbedPluginResources
}

// NavigationModuleInfo 导航所需要的模块信息
type NavigationModuleInfo struct {
	Name       string
	Title      string
	Path       string
	Navigation string
}
