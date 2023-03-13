package entry

import "time"

type ClusterConfigRuntime struct {
	Id          int
	NamespaceId int
	ConfigID    int
	ClusterId   int
	IsOnline    bool
	Operator    int
	CreateTime  time.Time
	UpdateTime  time.Time
}
