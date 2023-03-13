package entry

import "time"

type DiscoveryHistoryInfo struct {
	Discovery Discovery              `json:"discovery,omitempty"`
	Config    DiscoveryVersionConfig `json:"config,omitempty"`
}

type DiscoveryHistory struct {
	Id          int
	DiscoveryId int
	NamespaceId int
	OldValue    DiscoveryHistoryInfo
	NewValue    DiscoveryHistoryInfo
	OptType     OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}
