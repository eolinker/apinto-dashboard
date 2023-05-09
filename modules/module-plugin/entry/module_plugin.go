package entry

import (
	"time"
)

type ModulePlugin struct {
	Id         int       `gorm:"column:id" json:"id"`
	UUID       string    `gorm:"column:uuid" json:"uuid"`
	Name       string    `gorm:"column:name" json:"name"`
	Version    string    `gorm:"column:version" json:"version"`
	Group      string    `gorm:"column:group" json:"group"`
	Navigation string    `gorm:"column:navigation" json:"navigation"`
	CName      string    `gorm:"column:cname" json:"cname"`
	Resume     string    `gorm:"column:resume" json:"resume"`
	ICon       string    `gorm:"column:icon" json:"icon"`
	Type       int       `gorm:"column:type" json:"type"`
	Driver     string    `gorm:"column:driver" json:"driver"`
	Details    []byte    `gorm:"column:details" json:"details"`
	Operator   int       `gorm:"column:operator" json:"operator"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (*ModulePlugin) TableName() string {
	return "module_plugin"
}

func (c *ModulePlugin) IdValue() int {
	return c.Id
}

type EnablePlugin struct {
	UUID    string
	Name    string
	Driver  string
	Config  []byte
	Details []byte
}

// EnabledModule 启用的导航模块信息
type EnabledModule struct {
	Name            string `gorm:"column:name" json:"name"`
	Title           string `gorm:"column:cname" json:"cname"`
	Type            int    `gorm:"column:type" json:"type"`
	Navigation      string `gorm:"column:navigation" json:"navigation"`
	IsPluginVisible bool   `gorm:"column:is_plugin_visible" json:"is_plugin_visible"`
	Frontend        string `gorm:"column:frontend" json:"frontend"`
}

type PluginListItem struct {
	UUID     string `gorm:"column:uuid" json:"uuid"`
	Name     string `gorm:"column:name" json:"name"`
	CName    string `gorm:"column:cname" json:"cname"`
	Resume   string `gorm:"column:resume" json:"resume"`
	ICon     string `gorm:"column:icon" json:"icon"`
	Type     int    `gorm:"column:type" json:"type"`
	Group    string `gorm:"column:group" json:"group"`
	IsEnable int    `gorm:"is_enable" json:"is_enable"`
}
