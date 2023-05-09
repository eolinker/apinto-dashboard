package cluster_entry

import (
	"time"

	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
)

// Cluster 集群信息表
type Cluster struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"  json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:工作空间"  json:"namespace_id,omitempty"`
	Name        string    `gorm:"size:255;not null;column:name;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:集群名" json:"name,omitempty"`
	Title       string    `gorm:"column:title"`
	Desc        string    `gorm:"size:255;not null;column:desc;comment:集群名称"  json:"desc,omitempty"`
	Env         string    `gorm:"size:20;not null;column:env;comment:环境"  json:"env,omitempty"`
	Addr        string    `gorm:"size:255;not null;column:addr;comment:集群地址" json:"addr,omitempty"`
	UUID        string    `gorm:"size:255;not null;column:uuid;comment:集群ID" json:"uuid,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*Cluster) TableName() string {
	return "cluster"
}

func (c *Cluster) IdValue() int {
	return c.Id
}

// ClusterHistory 集群绑定的环境变量变更记录表
type ClusterHistory struct {
	Id          int
	ClusterId   int
	NamespaceId int
	OldValue    Cluster
	NewValue    Cluster
	OptType     history_entry.OptType
	Operator    int
	OptTime     time.Time
}
