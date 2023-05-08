package plugin_template_entry

import "time"

type PluginTemplatePublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	Target      int
	PluginTemplateVersionConfig
	OptType  int
	Operator int
	OptTime  time.Time
}
