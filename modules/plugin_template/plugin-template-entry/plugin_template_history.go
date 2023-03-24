package plugin_template_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type PluginTemplateHistory struct {
	Id               int
	PluginTemplateID int
	NamespaceId      int
	OldValue         PluginTemplateHistoryInfo
	NewValue         PluginTemplateHistoryInfo
	OptType          history_entry.OptType //1新增 2修改 3删除
	OptTime          time.Time
	Operator         int
}

type PluginTemplateHistoryInfo struct {
	PluginTemplate PluginTemplate              `json:"plugin_template"`
	Config         PluginTemplateVersionConfig `json:"config"`
}
