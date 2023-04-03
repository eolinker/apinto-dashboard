package namespace_entry

import "time"

// Namespace 工作空间
type Namespace struct {
	Id         int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Name       string    `gorm:"size:255;not null;column:name;comment:名称"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
}

func (Namespace) TableName() string {
	return "namespace"
}
