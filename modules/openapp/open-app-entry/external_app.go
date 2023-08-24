package open_app_entry

import (
	"time"
)

type ExternalApplication struct {
	Id         int       `gorm:"column:id;type:INT(11);AUTO_INCREMENT;primary_key;NOT NULL;comment:主键ID"`
	UUID       string    `gorm:"column:uuid;type:VARCHAR(36);uniqueIndex:unique_uuid;NOT NULL;comment:应用id，随机生成的长度32的字符串"`
	Namespace  int       `gorm:"column:namespace;type:INT(11);NOT NULL;comment:工作空间"`
	Name       string    `gorm:"column:name;type:VARCHAR(50);NOT NULL;comment:应用名称"`
	Token      string    `gorm:"column:token;type:VARCHAR(36);NOT NULL;comment:鉴权token"`
	Desc       string    `gorm:"column:desc;type:TEXT;comment:应用描述;"`
	Tags       string    `gorm:"column:tags;type:TEXT;comment:应用标签"`
	IsDisable  bool      `gorm:"column:is_disable;type:TINYINT(1); comment:是否禁用"`
	IsDelete   bool      `gorm:"column:is_delete;type:TINYINT(1); comment:是否删除"`
	Operator   int       `gorm:"column:operator;type:INT(11); comment:操作人"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime; comment:创建时间;not null;type:timestamp;not null;"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime;autoCreateTime;type:timestamp;DEFAULT:CURRENT_TIMESTAMP; comment:更新时间;"`
}

func (e *ExternalApplication) TableName() string {
	return "external_app"
}

func (e *ExternalApplication) IdValue() int {
	return e.Id
}
