package runtime_entry

import "time"

// Runtime 正在运行的版本表
type Runtime struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	ClusterID   int       `gorm:"type:int(11);size:11;not null;column:cluster;dbUniqueIndex:tag_cluster;uniqueIndex:tag_cluster;comment:集群ID"`
	Kind        string    `gorm:"size:20;not null;index:kind;column:kind;dbUniqueIndex:tag_cluster;uniqueIndex:tag_cluster;comment:类型"`
	TargetID    int       `gorm:"type:int(11);size:11;default:0;column:target;dbUniqueIndex:tag_cluster;uniqueIndex:tag_cluster;comment:根据type区分是哪个表的ID"`
	NamespaceID int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	Version     int       `gorm:"type:int(11);size:11;not null;column:version;comment:版本ID"`
	IsOnline    bool      `gorm:"type:tinyint(1);size:1;not null;default:1;column:is_online;comment:是否上线"`
	Disable     bool      `gorm:"type:tinyint(1);size:1;not null;column:disable;comment:禁用状态"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (r *Runtime) IdValue() int {
	return r.Id
}

func (*Runtime) TableName() string {
	return "runtime"
}
