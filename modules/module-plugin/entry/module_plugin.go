package entry

import (
	"time"
)

type ModulePlugin struct {
	Id         int       `gorm:"column:id" json:"id"`
	UUID       string    `gorm:"column:uuid" json:"uuid"`
	Name       string    `gorm:"column:name" json:"name"`
	Group      int       `gorm:"column:group" json:"group"`
	CName      string    `gorm:"column:cname" json:"cname"`
	Resume     string    `gorm:"column:resume" json:"resume"`
	ICon       string    `gorm:"column:icon" json:"icon"`
	Type       int       `gorm:"column:type" json:"type"`
	Driver     string    `gorm:"column:driver" json:"driver"`
	Details    []byte    `gorm:"column:details" json:"details"`
	Operator   int       `gorm:"column:operator" json:"operator"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (*ModulePlugin) TableName() string {
	return "module_plugin"
}

func (c *ModulePlugin) IdValue() int {
	return c.Id
}
