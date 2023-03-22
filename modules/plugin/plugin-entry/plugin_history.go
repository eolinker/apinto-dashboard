package plugin_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type PluginHistory struct {
	Id          int
	PluginID    int
	NamespaceId int
	OldValue    PluginHistoryInfo
	NewValue    PluginHistoryInfo
	OptType     history_entry.OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}

type PluginHistoryInfo struct {
	Plugin Plugin `json:"plugin"`
}
