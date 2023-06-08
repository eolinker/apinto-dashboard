package entry

type PluginResources struct {
	ID        int    `gorm:"type:int(11);size:11;not null;primary_key;column:id;comment:主键ID"`
	Uuid      string `gorm:"uniqueIndex:uuid;index;not null;size:36;column:uuid;comment:插件id"`
	Resources []byte `gorm:"type:mediumblob;serializer:json;column:resources"`
}

func (c *PluginResources) TableName() string {
	return "module_plugin_resources"
}

func (c *PluginResources) IdValue() int {
	return c.ID
}
