package entry

import "time"

type ServiceHistoryInfo struct {
	Service Service              `json:"service"`
	Config  ServiceVersionConfig `json:"config"`
}

type ServiceHistory struct {
	Id          int
	ServiceId   int
	NamespaceId int
	OldValue    ServiceHistoryInfo
	NewValue    ServiceHistoryInfo
	OptType     OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}
