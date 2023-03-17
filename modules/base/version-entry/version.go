package version_entry

import "time"

// Version 版本表
type Version struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceID int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	Target      int       `gorm:"type:int(11);size:11;not null;index:kind_target;column:target;comment:根据kind区分是哪个表的ID"`
	Kind        string    `gorm:"size:20;not null;column:kind;index:kind_target;comment:类型"`
	Data        []byte    `gorm:"type:text;not null;column:data;comment:配置信息json"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:操作人"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*Version) TableName() string {
	return "version"
}
func (s *Version) IdValue() int {
	return s.Id
}
