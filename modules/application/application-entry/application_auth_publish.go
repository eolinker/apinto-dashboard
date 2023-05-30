package application_entry

import "time"

type AppPublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	Target      int
	ApplicationVersionConfig
	OptType  int
	Operator int
	OptTime  time.Time
}
