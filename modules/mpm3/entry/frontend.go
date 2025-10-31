package entry

import "github.com/eolinker/apinto-dashboard/pm3"

type Frontend struct {
	Id      int           `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Plugin  int           `gorm:"type:int(11);size:11;not null;column:plugin;comment:所属插件id;index;"`
	Content pm3.PFrontend `gorm:"type:text;serializer:json;"`
}

func (c *Frontend) TableName() string {
	return "pm3_frontend"
}

func (c *Frontend) IdValue() int {
	return c.Id
}
