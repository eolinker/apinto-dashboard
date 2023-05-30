package dynamic_entry

import "time"

type Dynamic struct {
	Id          int       `gorm:"column:id"`
	NamespaceId int       `gorm:"column:namespace"`
	Name        string    `gorm:"column:name"`
	Title       string    `gorm:"column:title"`
	Skill       string    `gorm:"column:skill"`
	Driver      string    `gorm:"column:driver"`
	Description string    `gorm:"column:description"`
	Version     string    `gorm:"column:version"`
	Config      string    `gorm:"column:config"`
	Profession  string    `gorm:"column:profession"`
	Updater     int       `gorm:"column:updater"`
	CreateTime  time.Time `gorm:"column:create_time"`
	UpdateTime  time.Time `gorm:"column:update_time"`
}

func (d *Dynamic) TableName() string {
	return "dynamic_module"
}

func (d *Dynamic) IdValue() int {
	return d.Id
}
