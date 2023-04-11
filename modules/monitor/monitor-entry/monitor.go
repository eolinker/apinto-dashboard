package monitor_entry

import "time"

type MonitorPartition struct {
	Id         int       `json:"id,omitempty"`
	UUID       string    `json:"uuid,omitempty"`
	Namespace  int       `json:"namespace,omitempty"`
	Name       string    `json:"name,omitempty"`
	SourceType string    `json:"source_type,omitempty"`
	Config     []byte    `json:"config,omitempty"`
	Env        string    `json:"env,omitempty"`
	ClusterIDs string    `json:"cluster_ids,omitempty"`
	Operator   int       `json:"operator,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`
}

func (m *MonitorPartition) TableName() string {
	return "monitor"
}

func (m *MonitorPartition) IdValue() int {
	return m.Id
}
