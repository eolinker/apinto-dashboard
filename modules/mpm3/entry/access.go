package entry

type Access struct {
	Id     int      `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Plugin int      `gorm:"type:int(11);not null; column:plugin; comment:插件id; index"`
	Module string   `gorm:"size:255;not null;column:module;comment:模块标识"`
	Name   string   `gorm:"column:name;not null;size:255;comment:权限标识名"`
	CName  string   `gorm:"column:cname;not null;size:255;comment:权限中文名"`
	Depend []string `gorm:"column:depend;type:text;comment:依赖的权限;serializer:json"`
}

func (c *Access) TableName() string {
	return "pm3_access"
}

func (c *Access) IdValue() int {
	return c.Id
}
