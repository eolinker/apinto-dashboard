package upstream_entry

import (
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"time"
)

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
	OptType     history_entry.OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}
