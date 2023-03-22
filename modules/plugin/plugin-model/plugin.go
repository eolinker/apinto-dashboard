package plugin_model

import "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"

type Plugin struct {
	*plugin_entry.Plugin
	OperatorStr string
	IsDelete    bool
	IsBuilt     bool //是否内置
}

type PluginBasic struct {
	*plugin_entry.Plugin
}

type PluginInput struct {
	Name     string
	Extended string
	RelyName string
	Desc     string
}

type PluginBuilt struct {
	Extended string
	Name     string
	Rely     string
	Schema   string
	Sort     int
}

type ExtenderInfo struct {
	Id      string `json:"id"`
	Group   string `json:"group"`
	Project string `json:"project"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  string `json:"schema"`
}
