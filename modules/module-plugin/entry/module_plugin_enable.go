package entry

import "time"

type ModulePluginEnable struct {
	Id         int       `gorm:"id" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	Navigation string    `gorm:"navigation" json:"navigation"`
	IsEnable   int       `gorm:"is_enable" json:"is_enable"`
	Frontend   string    `gorm:"frontend" json:"frontend"`
	Config     []byte    `gorm:"config" json:"config"`
	Operator   int       `gorm:"operator" json:"operator"`
	UpdateTime time.Time `gorm:"update_time" json:"update_time"`
}

func (*ModulePluginEnable) TableName() string {
	return "module_plugin_enable"
}

func (p *ModulePluginEnable) IdValue() int {
	return p.Id
}
