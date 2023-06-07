package entry

import "gorm.io/gorm"

type PluginResources struct {
	Uuid      string      `gorm:"uniqueIndex:uuid;index;not null;size:36;column:uuid;comment:插件id"`
	Resources interface{} `gorm:"type:mediumblob;serializer:json;column:resources"`
	gorm.Model
}

func (c *PluginResources) TableName() string {
	return "module_plugin_resources"
}

func (c *PluginResources) IdValue() int {
	return int(c.ID)
}
