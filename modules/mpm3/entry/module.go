package entry

type Module struct {
	Id         int    `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Plugin     int    `gorm:"type:int(11);not null; column:plugin; comment:插件id; index"`
	Navigation string `gorm:"size:255;column:navigation;not null;comment:挂载导航id"`
	Name       string `gorm:"column:name;not null;size:255;comment:导航模块名"`
	CName      string `gorm:"column:cname;not null;size:255;comment:模块名中文名"`
	Router     string `gorm:"column:router;not null;size:255;comment:前端路由"`
}

func (*Module) TableName() string {
	return "pm3_module"
}

func (p *Module) IdValue() int {
	return p.Id
}
