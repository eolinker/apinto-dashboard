package audit_model

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type LogOperateType int

func (l LogOperateType) Handler(ginCtx *gin.Context) {
	ginCtx.Set("Operate", l)
}
func init() {
	for i := LogOperateTypeNone; i < LogOperateTypeALL; i++ {
		e := LogOperateType(i)
		logOperateTypeIndex[e.String()] = e
	}

	//for _, kind := range logKindOrder {
	//	logKindList = append(logKindList, LogKindsItem{
	//		Title: logKindTitles[kind],
	//		Name:  string(kind),
	//	})
	//}
}

const (
	LogOperateTypeNone LogOperateType = iota
	LogOperateTypeCreate
	LogOperateTypeEdit
	LogOperateTypeDelete
	LogOperateTypePublish
	LogOperateTypeALL
)

var (
	logOperateTypeNames = map[LogOperateType]string{
		LogOperateTypeNone:    "NONE",
		LogOperateTypeCreate:  "create",
		LogOperateTypeEdit:    "edit",
		LogOperateTypeDelete:  "delete",
		LogOperateTypePublish: "publish",
	}
	logOperateTypeTitles = map[LogOperateType]string{
		LogOperateTypeCreate:  "新建",
		LogOperateTypeEdit:    "编辑",
		LogOperateTypeDelete:  "删除",
		LogOperateTypePublish: "发布",
	}
	logOperateTypeIndex = map[string]LogOperateType{}
)

func (l LogOperateType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", l.String(), "\"")), nil
}

func (l LogOperateType) String() string {
	if l >= LogOperateTypeALL {
		return "unknown"
	}
	return logOperateTypeNames[l]
}

func GetLogOperateIndex(operateType string) LogOperateType {
	return logOperateTypeIndex[operateType]
}

func GetLogOperateTitle(operateIndex LogOperateType) string {
	return logOperateTypeTitles[operateIndex]
}

// LogKind 审计日志操作对象类型
type LogKind string

//const (
//	LogKindAPI             = "api"
//	LogKindService         = "service"
//	LogKindDiscovery       = "discovery"
//	LogKindApplication     = "application"
//	LogKindCluster         = "cluster"
//	LogKindGlobalVariable  = "global_variable"
//	LogKindClusterVariable = "cluster_variable"
//	LogKindStrategyTraffic = "strategy_traffic"
//	LogKindStrategyFuse    = "strategy_fuse"
//	LogKindStrategyVisit   = "strategy_visit"
//	LogKindStrategyCache   = "strategy_cache"
//	LogKindStrategyGrey    = "strategy_grey"
//	LogKindUser            = "user"
//	LogKindRole            = "role"
//	LogKindCommonGroup     = "common_group" //仅用于通用分组controller，用来标志是分组操作
//	LogKindAPIGroup        = "api_group"
//	LogKindExtAPP          = "ext_app"
//	LogKindMonPartition    = "monitor_partition" //监控分区
//	LogKindNoticeEmail     = "notice_email"      //通知渠道邮箱
//	LogKindNoticeWebhook   = "notice_webhook"    //通知渠道webhook
//	LogKindWarnStrategy    = "warn_strategy"     //告警策略
//	LogKindGlobalPlugin    = "global_plugin"
//	LogKindClusterPlugin   = "cluster_plugin"
//	LogKindPluginTemplate  = "plugin_template"
//	LogKindMiddleware      = "middleware"
//	LogKindNavigation      = "navigation"
//)
//
//var (
//	logKindTitles = map[LogKind]string{
//		LogKindAPI:             "API",
//		LogKindService:         "上游",
//		LogKindDiscovery:       "服务发现",
//		LogKindApplication:     "应用",
//		LogKindCluster:         "网关集群",
//		LogKindGlobalVariable:  "全局环境变量",
//		LogKindClusterVariable: "集群环境变量",
//		LogKindStrategyTraffic: "流量策略",
//		LogKindStrategyFuse:    "熔断策略",
//		LogKindStrategyVisit:   "访问策略",
//		LogKindStrategyCache:   "缓存策略",
//		LogKindStrategyGrey:    "灰度策略",
//		LogKindUser:            "用户",
//		LogKindRole:            "用户角色",
//		LogKindAPIGroup:        "API分组",
//		LogKindExtAPP:          "外部应用",
//		LogKindMonPartition:    "监控分区",
//		LogKindNoticeEmail:     "邮箱",
//		LogKindNoticeWebhook:   "webhook",
//		LogKindWarnStrategy:    "告警策略",
//		LogKindGlobalPlugin:    "全局插件",
//		LogKindClusterPlugin:   "集群插件",
//		LogKindPluginTemplate:  "插件模板",
//		LogKindMiddleware:      "拦截器分组",
//	}
//	logKindOrder = []LogKind{LogKindAPI, LogKindService, LogKindDiscovery, LogKindApplication, LogKindCluster, LogKindGlobalVariable, LogKindClusterVariable,
//		LogKindStrategyTraffic, LogKindStrategyFuse, LogKindStrategyVisit, LogKindStrategyCache, LogKindStrategyGrey,
//		LogKindAPIGroup, LogKindExtAPP}
//	logKindList []LogKindsItem
//)

type LogKindsItem struct {
	Title string `json:"title"`
	Name  string `json:"name"`
}

//
//func (l LogKind) MarshalJSON() ([]byte, error) {
//	return []byte(fmt.Sprint("\"", l.String(), "\"")), nil
//}
//
//func (l LogKind) String() string {
//	return logKindTitles[l]
//}
//
//func GetLogKinds() []LogKindsItem {
//	return logKindList
//}
//
//func GetLogKindTitle(kind LogKind) string {
//	return logKindTitles[kind]
//}
//
//func IsLogKindExist(kind string) bool {
//	_, has := logKindTitles[LogKind(kind)]
//	if has {
//		return true
//	}
//	return false
//}

type LogPublishType int

const (
	LogPublishTypeNone = iota
	LogPublishTypeOnline
	LogPublishTypeOffline
)

var (
	logPublishTypeTitles = map[LogPublishType]string{
		LogPublishTypeNone:    "NONE",
		LogPublishTypeOnline:  "上线",
		LogPublishTypeOffline: "下线",
	}
)

func GetPublishTypeTitle(publishType int) string {
	return logPublishTypeTitles[LogPublishType(publishType)]
}

type LogEnableOperateType int

const (
	LogEnableOperateTypeNone = iota
	LogEnableOperateTypeEnable
	LogEnableOperateTypeDisable
)

var (
	logEnableOperateTypeTitles = map[LogEnableOperateType]string{
		LogEnableOperateTypeNone:    "NONE",
		LogEnableOperateTypeEnable:  "启用",
		LogEnableOperateTypeDisable: "禁用",
	}
)

func GetEnableTypeTitle(enableType int) string {
	return logEnableOperateTypeTitles[LogEnableOperateType(enableType)]
}
