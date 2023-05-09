package strategy_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type StrategyHistoryInfo struct {
	Strategy Strategy           `json:"strategy,omitempty"`
	Config   StrategyConfigInfo `json:"config"`
}

type StrategyHistory struct {
	Id          int
	StrategyId  int
	NamespaceId int
	OldValue    StrategyHistoryInfo
	NewValue    StrategyHistoryInfo
	OptType     history_entry.OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}
