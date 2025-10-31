package entry

type Package struct {
	Id      int    `gorm:"column:id;type:int(11);autoIncrement:false;primary_key;comment:模块插件表的主键ID" json:"id"`
	Package []byte `gorm:"column:package;type:mediumblob;not null;comment:安装包" json:"package"`
}

func (*Package) TableName() string {
	return "module_plugin_package"
}

func (c *Package) IdValue() int {
	return c.Id
}
