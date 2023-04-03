package variable_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

// VariableHistory 集群绑定的环境变量变更记录表
type VariableHistory struct {
	Id          int
	ClusterId   int
	NamespaceId int
	VariableId  int
	OldValue    VariableValue
	NewValue    VariableValue
	OptType     history_entry.OptType
	Operator    int
	OptTime     time.Time
}

type VariableValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
