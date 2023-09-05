package dynamic_entry

import "time"

type Dynamic struct {
	Id          int       `gorm:"column:id;type:INT(11) UNSIGNED;AUTO_INCREMENT;NOT NULL;comment:主键ID"`
	NamespaceId int       `gorm:"column:namespace;not null;type:INT(11);uniqueIndex:unique;comment:工作空间"`
	Name        string    `gorm:"column:name;type:VARCHAR(255);NOT NULL;uniqueIndex:unique;comment:实例名"`
	Title       string    `gorm:"column:title;type:VARCHAR(255);comment:标题"`
	Driver      string    `gorm:"column:driver;type:VARCHAR(255);NOT NULL;comment:驱动"`
	Description string    `gorm:"column:description;type:VARCHAR(255);comment:描述"`
	Version     string    `gorm:"column:version;type:VARCHAR(32);NOT NULL;comment:版本"`
	Config      string    `gorm:"column:config;type:TEXT;NOT NULL;comment:配置"`
	Profession  string    `gorm:"column:profession;type:VARCHAR(255);NOT NULL;uniqueIndex:unique;comment:插件指定profession"`
	CreateTime  time.Time `gorm:"column:create_time;type:DATETIME;NOT NULL;comment:创建时间"`
	UpdateTime  time.Time `gorm:"column:update_time;type:DATETIME;NOT NULL;comment:更新时间"`
	Updater     int       `gorm:"column:updater;type:INT(11);NOT NULL;comment:更新人"`
	Skill       string    `gorm:"column:skill;type:VARCHAR(255);comment:模块提供能力"`
}

func (d *Dynamic) TableName() string {
	return "dynamic_module"
}

func (d *Dynamic) IdValue() int {
	return d.Id
}
