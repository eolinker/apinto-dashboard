package entry

import (
	"time"
)

type ModulePlugin struct {
	Id                  int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id"`
	UUID                string    `gorm:"size:36;not null;uniqueIndex:unique_uuid;column:uuid;comment:插件id" json:"uuid"`
	Name                string    `gorm:"size:255;not null;column:name;comment:插件名" json:"name"`
	Version             string    `gorm:"size:36;not null;column:version;comment:插件名" json:"version"`
	Group               string    `gorm:"size:255;not null;column:group;comment:插件分组id" json:"group"`
	Navigation          string    `gorm:"size:255;null;column:navigation;comment:导航id" json:"navigation"`
	CName               string    `gorm:"size:255;not null;column:cname;comment:昵称" json:"cname"`
	Resume              string    `gorm:"size:255;not null;column:resume;comment:简介" json:"resume"`
	ICon                string    `gorm:"size:255;not null;column:icon;comment:图标的文件名, 相对路径" json:"icon"`
	Driver              string    `gorm:"size:255;not null;column:driver;comment:插件类型" json:"driver"`
	IsCanDisable        bool      `gorm:"type:tinyint(1);size:1;not null;column:is_can_disable;comment:是否可停用" json:"is_can_disable"`
	IsCanUninstall      bool      `gorm:"type:tinyint(1);size:1;not null;column:is_can_uninstall;comment:是否可卸载" json:"is_can_uninstall"`
	IsInner             bool      `gorm:"type:tinyint(1);size:1;not null;column:is_inner;comment:是否为内置" json:"is_inner"`
	VisibleInNavigation bool      `gorm:"type:tinyint(1);size:1;not null;column:visible_in_navigation;comment:是否在导航显示" json:"visible_in_navigation"`
	VisibleInMarket     bool      `gorm:"type:tinyint(1);size:1;not null;column:visible_in_market;comment:是否在插件市场显示" json:"visible_in_market"`
	Details             []byte    `gorm:"type:mediumtext;not null;column:details;comment:插件define" json:"details"`
	Operator            int       `gorm:"type:int(11);size:11;column:operator;comment:操作人" json:"operator"`
	CreateTime          time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime          time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
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
	Name                string `gorm:"column:name" json:"name"`
	Title               string `gorm:"column:cname" json:"cname"`
	Navigation          string `gorm:"column:navigation" json:"navigation"`
	VisibleInNavigation bool   `gorm:"column:visible_in_navigation" json:"visible_in_navigation"`
	Frontend            string `gorm:"column:frontend" json:"frontend"`
}

type PluginListItem struct {
	UUID                string `gorm:"column:uuid" json:"uuid"`
	Name                string `gorm:"column:name" json:"name"`
	CName               string `gorm:"column:cname" json:"cname"`
	Resume              string `gorm:"column:resume" json:"resume"`
	ICon                string `gorm:"column:icon" json:"icon"`
	IsInner             bool   `gorm:"column:is_inner" json:"is_inner"`
	VisibleInNavigation bool   `gorm:"column:visible_in_navigation" json:"visible_in_navigation"`
	VisibleInMarket     bool   `gorm:"column:visible_in_market" json:"visible_in_market"`
	Group               string `gorm:"column:group" json:"group"`
	IsEnable            int    `gorm:"is_enable" json:"is_enable"`
}
