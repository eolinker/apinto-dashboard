package entry

type ModulePluginPackage struct {
	Id      int    `gorm:"id" json:"id"`
	Package []byte `gorm:"package" json:"package"`
}

func (*ModulePluginPackage) TableName() string {
	return "module_plugin_package"
}

func (c *ModulePluginPackage) IdValue() int {
	return c.Id
}
