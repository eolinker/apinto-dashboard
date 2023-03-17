package discovery_entry

import "time"

type DiscoveryRuntime struct {
	Id          int
	NamespaceId int
	DiscoveryID int
	ClusterId   int
	VersionID   int
	IsOnline    bool
	Operator    int
	CreateTime  time.Time
	UpdateTime  time.Time
}
