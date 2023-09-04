package entry

import "time"

type ModulePluginEnable struct {
	Id         int       `gorm:"type:int(11);primary_key;column:id;autoIncrement:false;comment:模块插件的主键id" json:"id"`
	Name       string    `gorm:"column:name;not null;size:255;comment:模块名，可以改，默认为模块插件的name" json:"name"`
	Navigation string    `gorm:"column:navigation;not null;size:255;comment:导航id" json:"navigation"`
	IsEnable   int       `gorm:"column:is_enable;not null;type:tinyint(1);comment:是否启用 1未启用 2启用" json:"is_enable"`
	Frontend   string    `gorm:"column:frontend;size:255;comment:前端路由" json:"frontend"`
	Config     []byte    `gorm:"column:config;type:text;not null;comment:启用配置" json:"config"`
	Operator   int       `gorm:"column:operator;type:int(11);comment:更新人/操作人" json:"operator"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;autoCreateTime;autoUpdateTime;comment:修改时间" json:"update_time"`
}

func (*ModulePluginEnable) TableName() string {
	return "module_plugin_enable"
}

func (p *ModulePluginEnable) IdValue() int {
	return p.Id
}
