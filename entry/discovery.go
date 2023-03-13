package entry

import "time"

type Discovery struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:工作空间" json:"namespace_id,omitempty"`
	Name        string    `gorm:"size:255;not null;column:name;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:名称" json:"name,omitempty"`
	UUID        string    `gorm:"size:36;not null;column:uuid" json:"uuid,omitempty"` //TODO 这里应该才是唯一键，后续改
	Driver      string    `gorm:"size:100;default:null;column:driver;comment:驱动" json:"driver,omitempty"`
	Desc        string    `gorm:"size:255;default:null;column:desc;comment:描述" json:"desc,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*Discovery) TableName() string {
	return "discovery"
}

func (s *Discovery) IdValue() int {
	return s.Id
}
