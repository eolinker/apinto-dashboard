package entry

import "time"

// ClusterNode 集群绑定的节点信息表
type ClusterNode struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;uniqueIndex:cluster_name;column:cluster;comment:集群ID"`
	Name        string    `gorm:"size:255;not null;column:name;uniqueIndex:cluster_name;comment:节点名称/ID"`
	AdminAddr   string    `gorm:"size:255;default:null;column:admin_addr;comment:管理地址"`
	ServiceAddr string    `gorm:"size:255;not null;column:service_addr;comment:服务地址"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
}

func (*ClusterNode) TableName() string {
	return "cluster_node"
}

func (c *ClusterNode) GetKey() string {
	return c.Name
}

func (c *ClusterNode) Values() []string {
	return []string{c.AdminAddr}
}

func (c *ClusterNode) IdValue() int {
	return c.Id
}
