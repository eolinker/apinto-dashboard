package model

import (
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
)

type ModulePluginItem struct {
	*entry.ModulePlugin
	IsEnable bool
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
	Header     []ExtendParams
	Query      []ExtendParams
	Initialize []ExtendParams
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
	Name  string
	Value string
}

type ExtendParamsRender struct {
	Name        string
	Value       string
	Title       string
	Placeholder string
	Desc        string
}
