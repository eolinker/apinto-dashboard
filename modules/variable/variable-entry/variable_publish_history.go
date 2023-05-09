package variable_entry

import "time"

type VariablePublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	VariablePublishHistoryInfo
	OptType  int
	Operator int
	OptTime  time.Time
}

type VariablePublishHistoryInfo struct {
	VariableToPublish []*VariableToPublish `json:"variable_to_publish"` //发布记录
}
