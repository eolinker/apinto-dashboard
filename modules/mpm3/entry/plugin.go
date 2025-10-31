package entry

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"time"
)

type EnablePlugin struct {
	UUID    string            `gorm:"size:36;not null;uniqueIndex:unique_uuid;column:uuid;comment:插件id" json:"uuid"`
	Name    string            `gorm:"size:255;not null;column:name;comment:插件名" json:"name"`
	Driver  string            `gorm:"size:255;not null;column:driver;comment:插件类型" json:"driver"`
	Config  []byte            `gorm:"column:config;type:text;comment:启用配置;serializer:json" json:"config"`
	Details *pm3.PluginDefine `gorm:"type:mediumtext;not null;column:details;comment:插件define;serializer:json" json:"details"`
}
type Plugin struct {
	Id           int               `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id"`
	UUID         string            `gorm:"size:36;not null;uniqueIndex:unique_uuid;column:uuid;comment:插件id" json:"uuid"`
	Name         string            `gorm:"size:255;not null;column:name;comment:插件名" json:"name"`
	Version      string            `gorm:"size:36;not null;column:version;comment:插件版本" json:"version"`
	Group        string            `gorm:"size:255;not null;column:group;comment:插件分组id" json:"group"`
	CName        string            `gorm:"size:255;not null;column:cname;comment:昵称" json:"cname"`
	Resume       string            `gorm:"size:255;not null;column:resume;comment:简介" json:"resume"`
	ICon         string            `gorm:"size:255;not null;column:icon;comment:图标的文件名, 相对路径" json:"icon"`
	Driver       string            `gorm:"size:255;not null;column:driver;comment:插件类型" json:"driver"`
	IsCanDisable bool              `gorm:"type:tinyint(1);size:1;not null;column:is_can_disable;comment:是否可停用" json:"is_can_disable"`
	IsInner      bool              `gorm:"type:tinyint(1);size:1;not null;column:is_inner;comment:是否为内置" json:"is_inner"`
	Details      *pm3.PluginDefine `gorm:"type:mediumtext;not null;column:details;comment:插件define;serializer:json" json:"details"`
	Operator     int               `gorm:"type:int(11);size:11;column:operator;comment:操作人" json:"operator"`
	Hash         string            `gorm:"size:36;not null;column:hash;comment:hash" json:"hash"`
	CreateTime   time.Time         `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime   time.Time         `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*Plugin) TableName() string {
	return "pm3_plugins"
}

func (c *Plugin) IdValue() int {
	return c.Id
}

type PluginListItem struct {
	UUID    string `gorm:"column:uuid" json:"uuid"`
	Name    string `gorm:"column:name" json:"name"`
	CName   string `gorm:"column:cname" json:"cname"`
	Resume  string `gorm:"column:resume" json:"resume"`
	ICon    string `gorm:"column:icon" json:"icon"`
	IsInner bool   `gorm:"column:is_inner" json:"is_inner"`

	Group    string `gorm:"column:group" json:"group"`
	IsEnable int    `gorm:"is_enable" json:"is_enable"`
}
