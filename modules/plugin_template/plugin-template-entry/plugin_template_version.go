package plugin_template_entry

import "time"

type PluginTemplateVersion struct {
	Id               int
	PluginTemplateId int
	NamespaceID      int
	PluginTemplateVersionConfig
	Operator   int
	CreateTime time.Time
}

func (p *PluginTemplateVersion) SetVersionId(id int) {
	p.Id = id
}

type PluginTemplateVersionConfig struct {
	Plugins []*PluginTemplateVersionConfigDetail `json:"plugins"`
}

type PluginTemplateVersionConfigDetail struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	Config  string `json:"config,omitempty"`
	Disable bool   `json:"disable,omitempty"`
}
