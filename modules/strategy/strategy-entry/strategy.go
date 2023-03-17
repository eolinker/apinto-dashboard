package strategy_entry

import "time"

// Strategy 策略表
type Strategy struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	UUID        string    `gorm:"size:36;not null;column:uuid;dbUniqueIndex:uuid;uniqueIndex:uuid;comment:uuid"  json:"uuid,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace_id,omitempty"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;column:cluster;index:cluster_priority_type;comment:集群ID" json:"cluster_id,omitempty"`
	Type        string    `gorm:"size:50;default null;column:type;index:cluster_priority_type;comment:策略类型" json:"type,omitempty"`
	Name        string    `gorm:"size:50;not null;column:name;comment:策略名" json:"name,omitempty"`
	Desc        string    `gorm:"size:255;default:null;column:desc;comment:描述" json:"desc,omitempty"`
	Priority    int       `gorm:"type:int(11);size:11;not null;column:priority;index:cluster_priority_type;comment:优先级" json:"priority,omitempty"`
	IsStop      bool      `gorm:"type:tinyint(1);size:1;column:is_stop;comment:是否停用" json:"is_stop,omitempty"` //true 停用  false启用
	IsDelete    bool      `gorm:"type:tinyint(1);size:1;column:is_delete;comment:是否已经删除(软删除)" json:"is_delete,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (s *Strategy) TableName() string {
	return "strategy"
}

func (s *Strategy) IdValue() int {
	return s.Id
}
