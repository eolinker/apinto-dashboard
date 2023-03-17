package variable_entry

import "time"

// Variables 全局环境变量信息表
type Variables struct {
	Id         int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Namespace  int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_key_index;uniqueIndex:namespace_key_index;comment:工作空间"`
	Key        string    `gorm:"size:255;not null;column:key;dbUniqueIndex:namespace_key_index;uniqueIndex:namespace_key_index;comment:key"`
	Desc       string    `gorm:"size:255;not null;column:desc;comment:描述"`
	Operator   int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (*Variables) TableName() string {
	return "variables"
}

func (v *Variables) IdValue() int {
	return v.Id
}
