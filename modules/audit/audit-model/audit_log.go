package audit_model

import "github.com/eolinker/apinto-dashboard/enum"

type LogListItem struct {
	ID          int                 `json:"id"`
	Operator    *OperatorInfo       `json:"operator"`
	OperateType enum.LogOperateType `json:"operate_type"`
	Kind        enum.LogKind        `json:"kind"`
	Time        string              `json:"time"`
	IP          string              `json:"ip"`
}

type LogDetailArg struct {
	Attr  string `json:"attr"`
	Value string `json:"value"`
}

// OperatorInfo 操作人信息
type OperatorInfo struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// LogObjectInfo 日志目标对象信息
type LogObjectInfo struct {
	Uuid          string `json:"uuid,omitempty"`
	Name          string `json:"name,omitempty"`
	ClusterId     int    `json:"cluster_id,omitempty"`
	ClusterName   string `json:"cluster_name,omitempty"`
	EnableOperate int    `json:"enable_operate,omitempty"` //启用禁用操作类型 1为启用，2为禁用，针对有启用禁用的模块
	PublishType   int    `json:"publish_type,omitempty"`   //发布类型 1为上线，2为下线 针对api，service这类有上下线区分的
}
