package strategy_entry

import "time"

type StrategyPublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	Publish     []*StrategyPublishConfigInfo //此次发布的数据
	OptType     int
	Operator    int
	CreateTime  time.Time
}

//type StrategyPublishHistoryInfo struct {
//	Publish []*StrategyPublishConfigInfo `json:"publish"`
//}
