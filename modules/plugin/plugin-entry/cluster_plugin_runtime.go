package plugin_entry

import "time"

// ClusterPluginRuntime 集群当前插件配置的版本
type ClusterPluginRuntime struct {
	Id          int       `json:"id"`
	ClusterId   int       `json:"cluster_id"`
	NamespaceId int       `json:"namespace_id"`
	VersionId   int       `json:"version_id"`
	IsOnline    bool      `json:"is_online"`
	Operator    int       `json:"operator"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}
