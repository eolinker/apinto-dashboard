package dynamic_entry

import "time"

type DynamicPublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	Publish     *DynamicPublishConfig //此次发布的数据
	OptType     int
	Operator    int
	CreateTime  time.Time
}
