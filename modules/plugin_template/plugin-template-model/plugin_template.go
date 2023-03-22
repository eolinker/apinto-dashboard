package plugin_template_model

import (
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"time"
)

type PluginTemplate struct {
	*plugin_template_entry.PluginTemplate
	OperatorStr string
	IsDelete    bool
}

type PluginTemplateDetail struct {
	*PluginTemplate
	Plugins []*PluginInfo
}

type PluginTemplateBasicInfo struct {
	*plugin_template_entry.PluginTemplate
}

type PluginTemplateVersion plugin_template_entry.PluginTemplateVersion

type PluginInfo struct {
	Name    string
	Config  string
	Disable bool
}

type PluginTemplateOnlineItem struct {
	ClusterName string
	ClusterEnv  string
	Status      int
	Disable     bool
	Operator    string
	UpdateTime  time.Time
}
