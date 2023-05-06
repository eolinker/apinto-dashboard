package api_entry

import "time"

type ApiPublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	Target      int
	APIVersionConfig
	OptType  int
	Operator int
	OptTime  time.Time
}
