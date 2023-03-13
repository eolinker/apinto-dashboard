package entry

import "time"

type ApplicationAuthPublish struct {
	Id              int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceId     int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace,omitempty"`
	Cluster         int       `gorm:"type:int(11);size:11;not null;dbUniqueIndex:cluster_app_auth;uniqueIndex:cluster_app_auth;column:cluster;comment:集群ID"`
	Application     int       `gorm:"type:int(11);size:11;not null;dbUniqueIndex:cluster_app_auth;uniqueIndex:cluster_app_auth;column:application;comment:application表ID"`
	ApplicationAuth int       `gorm:"type:int(11);size:11;not null;dbUniqueIndex:cluster_app_auth;uniqueIndex:cluster_app_auth;column:application_auth;comment:application_auth表ID"`
	Operator        int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime      time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
}

func (ApplicationAuthPublish) TableName() string {
	return "application_auth_publish"
}
