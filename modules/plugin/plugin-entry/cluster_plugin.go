package plugin_entry

import (
	"time"
)

type ClusterPlugin struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace,omitempty"`
	ClusterId   int       `gorm:"type:int(11);size:11;not null;column:cluster;dbUniqueIndex:cluster_plugin;comment:集群ID" json:"cluster,omitempty"`
	PluginName  string    `gorm:"size:255;not null;column:plugin_name;dbUniqueIndex:cluster_plugin;comment:插件名称" json:"plugin_name,omitempty"`
	Status      int       `gorm:"type:tinyint(2);not null;column:status;comment:插件状态" json:"status,omitempty"` //1禁用 2启用 3全局启用
	Config      string    `gorm:"type:text;column:config;comment:插件配置信息" json:"config,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*ClusterPlugin) TableName() string {
	return "cluster_plugin"
}

func (c *ClusterPlugin) IdValue() int {
	return c.Id
}
