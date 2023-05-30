package history_entry

import "time"

type HistoryKind = string
type OptType int

const (
	OptAdd OptType = iota + 1
	OptEdit
	OptDelete
)
const (
	HistoryKindAPI             HistoryKind = "api"
	HistoryKindApplication     HistoryKind = "application"
	HistoryKindApplicationAuth HistoryKind = "application_auth"
	HistoryKindCluster         HistoryKind = "cluster"
	HistoryKindDiscovery       HistoryKind = "discovery"
	HistoryKindService         HistoryKind = "service"
	HistoryKindStrategy        HistoryKind = "strategy"
	HistoryKindVariable        HistoryKind = "variable"
	HistoryKindVariableGlobal  HistoryKind = "variable_global"
	HistoryKindRole            HistoryKind = "role"
	HistoryKindPluginTemplate  HistoryKind = "plugin_template"
	HistoryKindPlugin          HistoryKind = "plugin"
	HistoryKindClusterPlugin   HistoryKind = "cluster_plugin"
	HistoryKindDynamicModule   HistoryKind = "dynamic_module"
)

// History 变更记录表
type History struct {
	Id          int         `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceID int         `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	Kind        HistoryKind `gorm:"size:50;not null;column:kind;comment:kind"`
	TargetID    int         `gorm:"type:int(11);size:11;not null;column:target;comment:根据kind区分"`
	OldValue    string      `gorm:"type:text;column:old_value;comment:oldValue"`
	NewValue    string      `gorm:"type:text;column:new_value;comment:newValue"`
	OptType     OptType     `gorm:"type:tinyint(4);size:4;default:0;column:opt_type;comment:1新增 2修改 3删除"` //1新增 2修改 3删除
	Operator    int         `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人"`
	OptTime     time.Time   `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:opt_time;comment:创建时间"`
}

func (*History) TableName() string {
	return "history"
}
