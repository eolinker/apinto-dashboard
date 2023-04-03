package api_entry

import "time"

type APIRuntime struct {
	Id          int
	NamespaceId int
	ApiID       int
	ClusterID   int
	VersionID   int
	IsOnline    bool
	Disable     bool
	Operator    int
	CreateTime  time.Time
	UpdateTime  time.Time
}
