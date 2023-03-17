package plugin_template_entry

import "time"

type PluginTemplateRuntime struct {
	Id               int
	NamespaceId      int
	PluginTemplateID int
	ClusterID        int
	VersionID        int
	IsOnline         bool
	Disable          bool
	Operator         int
	CreateTime       time.Time
	UpdateTime       time.Time
}
