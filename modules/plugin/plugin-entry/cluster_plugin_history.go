package plugin_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

// ClusterPluginHistory 集群绑定的插件变更记录表
type ClusterPluginHistory struct {
	Id              int
	NamespaceId     int
	ClusterPluginID int
	OldValue        ClusterPluginHistoryValue
	NewValue        ClusterPluginHistoryValue
	OptType         history_entry.OptType
	Operator        int
	OptTime         time.Time
}

type ClusterPluginHistoryValue struct {
	PluginName string `json:"plugin_name"`
	Status     int    `json:"status"`
	Config     string `json:"config"`
}
