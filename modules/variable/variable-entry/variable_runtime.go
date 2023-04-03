package variable_entry

import "time"

// VariableRuntime 集群当前版本
type VariableRuntime struct {
	Id          int       `json:"id"`
	ClusterId   int       `json:"cluster_id"`
	NamespaceId int       `json:"namespace_id"`
	VersionId   int       `json:"version_id"`
	IsOnline    bool      `json:"is_online"`
	Operator    int       `json:"operator"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}
