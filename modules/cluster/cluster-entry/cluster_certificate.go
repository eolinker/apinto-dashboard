package cluster_entry

import "time"

// ClusterCertificate 集群绑定的证书信息表
type ClusterCertificate struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:证书ID"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;index:cluster;column:cluster;comment:集群ID"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	Operator    int       `gorm:"type:int(11);size:11;default:null;column:operator;comment:更新人/操作人"`
	Key         string    `gorm:"type:text;not null;column:key"`
	Pem         string    `gorm:"type:text;not null;column:pem"`
	UUID        string    `gorm:"size:36;not null;column:uuid;dbUniqueIndex:uuid;uniqueIndex:uuid;comment:uuid"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (*ClusterCertificate) TableName() string {
	return "cluster_certificate"
}

func (c *ClusterCertificate) IdValue() int {
	return c.Id
}
