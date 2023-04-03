package variable_entry

import "time"

// ClusterVariable 集群绑定的环境变量信息表
type ClusterVariable struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;column:cluster;dbUniqueIndex:cluster_variable_index;uniqueIndex:cluster_variable_index;comment:集群ID" json:"cluster_id"`
	VariableId  int       `gorm:"type:int(11);size:11;not null;column:variable;dbUniqueIndex:cluster_variable_index;uniqueIndex:cluster_variable_index;comment:全局环境变量ID" json:"variable_id"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace_id"`
	Status      int       `gorm:"type:int(11);size:11;not null;default:0;column:status;comment:状态 1删除" json:"status"`
	Key         string    `gorm:"size:255;not null;column:key" json:"key"`
	Value       string    `gorm:"size:255;not null;column:value" json:"value"`
	Operator    int       `gorm:"type:int(11);size:11;default:null;column:operator;comment:操作人"` //操作人
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (*ClusterVariable) TableName() string {
	return "variable_cluster"
}

func (c *ClusterVariable) GetKey() string {
	return c.Key
}

func (c *ClusterVariable) Values() []string {
	return []string{c.Value}
}

func (c *ClusterVariable) IdValue() int {
	return c.Id
}
