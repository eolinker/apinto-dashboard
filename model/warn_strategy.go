package model

import "time"

type QueryWarnStrategyParam struct {
	PartitionId  int
	StartTime    time.Time
	EndTime      time.Time
	StrategyName string
	Dimension    []string
	Status       int
	PageNum      int
	PageSize     int
}

type WarnStrategy struct {
	PartitionId        int                 `json:"partition_id"`
	NamespaceId        int                 `json:"namespace_id"`
	Uuid               string              `json:"uuid"`
	Title              string              `json:"title"`
	Desc               string              `json:"desc"`
	IsEnable           bool                `json:"is_enable"`
	Dimension          string              `json:"dimension"` // api/service/cluster/partition
	Quota              QuotaType           `json:"quota"`
	Every              int                 `json:"every"`
	WarnStrategyConfig *WarnStrategyConfig `json:"warn_strategy_config"`
	PartitionUUID      string
	Operator           string
	CreateTime         time.Time
	UpdateTime         time.Time
}

const (
	DimensionTypeApi       = "api"
	DimensionTypeService   = "service"
	DimensionTypeCluster   = "cluster"
	DimensionTypePartition = "partition"
)

type WarnStrategyConfig struct {
	Target     WarnStrategyConfigTarget  `json:"target"`          //告警目标
	Rule       []*WarnStrategyConfigRule `json:"rule"`            //告警规则
	Continuity int                       `json:"continuity"`      //连续告警每N分钟1次
	HourMax    int                       `json:"hour_max"`        //每小时最多N次
	Users      []int                     `json:"users,omitempty"` //用户ID数组
}

type QuotaType string

var QuotaTypeReqFailCount QuotaType = "request_fail_count" //请求失败数
var QuotaTypeReqFailRate QuotaType = "request_fail_rate"   //请求失败率
var QuotaTypeReqStatus4xx QuotaType = "request_status_4xx" //请求4xx状态码数
var QuotaTypeReqStatus5xx QuotaType = "request_status_5xx" //请求5xx状态码数
var QuotaTypeProxyFailCount QuotaType = "proxy_fail_count" //转发失败数
var QuotaTypeProxyFailRate QuotaType = "proxy_fail_rate"   //转发失败率
var QuotaTypeProxyStatus4xx QuotaType = "proxy_status_4xx" //转发4xx状态码数
var QuotaTypeProxyStatus5xx QuotaType = "proxy_status_5xx" //转发5xx状态码数
var QuotaTypeReqMessage QuotaType = "request_message"      //请求报文量
var QuotaTypeRespMessage QuotaType = "response_message"    //响应报文量
var QuotaTypeAvgResp QuotaType = "avg_resp"                //平均响应时间

var QuotaRuleMap = map[QuotaType]string{QuotaTypeReqFailCount: "请求失败数", QuotaTypeReqFailRate: "请求失败率", QuotaTypeReqStatus4xx: "请求4xx状态码数",
	QuotaTypeReqStatus5xx: "请求5xx状态码数", QuotaTypeProxyFailCount: "转发失败数", QuotaTypeProxyFailRate: "转发失败率", QuotaTypeProxyStatus4xx: "转发4xx状态码数",
	QuotaTypeProxyStatus5xx: "转发5xx状态码数", QuotaTypeReqMessage: "请求报文量", QuotaTypeRespMessage: "响应报文量", QuotaTypeAvgResp: "平均响应时间"}

type WarnStrategyConfigTarget struct {
	Rule   string   `json:"rule"`   // 不限：unlimited 包含：contain 不包含：not_contain
	Values []string `json:"values"` // dimension为API则为API数组 以此类推
}

const (
	RuleTypeUnlimited  = "unlimited"
	RuleTypeContain    = "contain"
	RuleTypeNotContain = "not_contain"
)

type WarnStrategyConfigRule struct {
	ChannelUuids []string                           `json:"channel_uuids"` //渠道uuid数组 []string
	Condition    []*WarnStrategyConfigRuleCondition `json:"condition"`     //条件组
}
type WarnStrategyConfigRuleCondition struct {
	Compare string  `json:"compare"` //比较关系
	Unit    string  `json:"unit"`    //比较单位 num  %  ms  kb
	Value   float64 `json:"value"`   //阈值
}

var CompareValue = map[string]string{">": ">", ">=": ">=", "<": "<", "<=": "<=", "==": "==", "!=": "!=", "ring_ratio_add": "环比上个时间增率", "ring_ratio_reduce": "环比上个时间减率", "year_basis_add": "同比昨天同时间增率", "year_basis_reduce": "同比昨天同时间减率"}

type NoticeChannelWarn struct {
	Url       string
	Name      string
	Every     int
	Quota     QuotaType
	Condition []*MsgCondition
}

type MsgCondition struct {
	RealValue float64 //influx中的实际值
	Compare   string  `json:"compare"` //比较关系
	Unit      string  `json:"unit"`    //比较单位 num  %  ms  kb
	Value     float64 `json:"value"`   //阈值
}

type MsgWebhook struct {
	Title   string `json:"标题"`
	Name    string `json:"告警目标"`
	Url     string `json:"url,omitempty"`
	Content string `json:"告警内容"`
	Time    string `json:"告警时间"`
}

const (
	// WarnStatusNotTrigger 未触发触发告警
	WarnStatusNotTrigger = "not_trigger"
	// WarnStatusTrigger 触发告警（未发送）
	WarnStatusTrigger = "trigger"
	// WarnStatusSendTrigger 触发告警（已发送）
	WarnStatusSendTrigger = "send_trigger"
)

type SendMsgError struct {
	UUID              string `json:"uuid"`
	NoticeChannelUUID string `json:"notice_channel_uuid"`
	Msg               string `json:"msg"`
}
