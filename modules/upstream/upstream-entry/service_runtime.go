package upstream_entry

import "time"

type ServiceRuntime struct {
	Id          int
	NamespaceId int
	ServiceId   int
	ClusterId   int
	VersionId   int
	IsOnline    bool
	Operator    int
	CreateTime  time.Time
	UpdateTime  time.Time
}
