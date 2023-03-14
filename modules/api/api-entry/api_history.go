package api_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type ApiHistory struct {
	Id          int
	ApiId       int
	NamespaceId int
	OldValue    ApiHistoryInfo
	NewValue    ApiHistoryInfo
	OptType     history_entry.OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}

type ApiHistoryInfo struct {
	Api    API              `json:"api"`
	Config APIVersionConfig `json:"config"`
}
