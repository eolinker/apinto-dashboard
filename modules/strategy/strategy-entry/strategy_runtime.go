package strategy_entry

import "time"

type StrategyRuntime struct {
	Id          int       `json:"id"`
	ClusterId   int       `json:"cluster_id"`
	NamespaceId int       `json:"namespace_id"`
	VersionId   int       `json:"version_id"`
	IsOnline    bool      `json:"is_online"`
	Operator    int       `json:"operator"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}
