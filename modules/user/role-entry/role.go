package role_entry

import "time"

// Role 角色表
type Role struct {
	ID         int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id"`
	Title      string    `gorm:"size:16;not null;column:title;dbUniqueIndex:unique_title;uniqueIndex:unique_title;comment:角色名" json:"title"`
	Uuid       string    `gorm:"size:36;not null;column:uuid;dbUniqueIndex:unique_uuid;uniqueIndex:unique_uuid;comment:uuid" json:"uuid"`
	Desc       string    `gorm:"size:255;default null;column:desc;comment:描述" json:"desc"`
	Module     string    `gorm:"size:512;not null;column:module;comment:模块目录（多级模块用‘/’切割）" json:"module"`
	Type       int       `gorm:"type:tinyint(1);size:1;default:1;column:type;comment:类型(0=内置,1=自定义)" json:"type"`
	Operator   int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (r *Role) IdValue() int {
	return r.ID
}

func (*Role) TableName() string {
	return "role"
}
