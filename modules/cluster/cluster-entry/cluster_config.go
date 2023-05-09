package cluster_entry

import "time"

// ClusterConfig 集群配置信息表
type ClusterConfig struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	Type        string    `gorm:"size:50;default:null;column:type;dbUniqueIndex:cluster_type;uniqueIndex:cluster_type"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;column:cluster;dbUniqueIndex:cluster_type;uniqueIndex:cluster_type;comment:集群ID"`
	IsEnable    bool      `gorm:"type:tinyint(1);size:1;not null;column:is_enable;comment:是否启用"`
	Data        []byte    `gorm:"type:text;not null;column:data;comment:配置数据"`
	Operator    int       `gorm:"type:int(11);size:11;default:null;column:operator;comment:更新人/操作人"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (*ClusterConfig) TableName() string {
	return "cluster_config"
}

func (c *ClusterConfig) IdValue() int {
	return c.Id
}
