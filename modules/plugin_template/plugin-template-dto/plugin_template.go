package plugin_template_dto

import "github.com/eolinker/apinto-dashboard/enum"

type PluginTemplate struct {
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Operator   string `json:"operator"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	IsDelete   bool   `json:"is_delete"`
}

type PluginTemplateInput struct {
	Uuid    string        `json:"uuid"`
	Name    string        `json:"name"`
	Desc    string        `json:"desc"`
	Plugins []*PluginInfo `json:"plugins"`
}

type PluginInfo struct {
	Name    string      `json:"name"`
	Config  interface{} `json:"config"`
	Disable bool        `json:"disable"`
}

type PluginTemplateOutput struct {
	Name    string        `json:"name"`
	Desc    string        `json:"desc"`
	Plugins []*PluginInfo `json:"plugins"`
}

type OnlineOut struct {
	Name       string            `json:"name"`
	Status     enum.OnlineStatus `json:"status"`
	Env        string            `json:"env"`
	Disable    bool              `json:"disable"`
	Operator   string            `json:"operator"`
	UpdateTime string            `json:"update_time"`
}

type OnlineInput struct {
	ClusterName string `json:"cluster_name"`
}
