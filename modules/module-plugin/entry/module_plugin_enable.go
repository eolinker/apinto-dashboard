package entry

import "time"

type ModulePluginEnable struct {
	Id              int       `gorm:"id" json:"id"`
	Name            string    `gorm:"name" json:"name"`
	Navigation      string    `gorm:"navigation" json:"navigation"`
	IsEnable        int       `gorm:"is_enable" json:"is_enable"`
	IsCanDisable    bool      `gorm:"is_can_disable" json:"is_can_disable"`
	IsCanUninstall  bool      `gorm:"is_can_uninstall" json:"is_can_uninstall"`
	IsShowServer    bool      `gorm:"is_show_server" json:"is_show_server"`
	IsPluginVisible bool      `gorm:"is_plugin_visible" json:"is_plugin_visible"`
	Frontend        string    `gorm:"frontend" json:"frontend"`
	Config          []byte    `gorm:"config" json:"config"`
	Operator        int       `gorm:"operator" json:"operator"`
	UpdateTime      time.Time `gorm:"update_time" json:"update_time"`
}

func (*ModulePluginEnable) TableName() string {
	return "module_plugin_enable"
}

func (p *ModulePluginEnable) IdValue() int {
	return p.Id
}
