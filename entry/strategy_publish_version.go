package entry

import "time"

type StrategyPublishVersion struct {
	Id          int                          `json:"id"`
	ClusterId   int                          `json:"cluster_id"`
	NamespaceId int                          `json:"namespace_id"`
	Publish     []*StrategyPublishConfigInfo `json:"publish"`
	Operator    int                          `json:"operator"`
	CreateTime  time.Time                    `json:"create_time"`
}

func (v *StrategyPublishVersion) SetVersionId(id int) {
	v.Id = id
}

type StrategyPublishConfigInfo struct {
	Strategy        Strategy        `json:"strategy"`
	StrategyVersion StrategyVersion `json:"strategy_version"`
	Status          int             `json:"status"`
}

//type StrategyVersionPublishConfig struct {
//	Publish []*StrategyPublishConfigInfo `json:"publish"`
//}
