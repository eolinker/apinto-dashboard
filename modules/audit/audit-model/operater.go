package audit_model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type LogOperateType int

func (l LogOperateType) Handler(ginCtx *gin.Context) {
	ginCtx.Set("Operate", int(l))
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
		LogOperateTypeNone:    "none",
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
	l, has := logOperateTypeIndex[strings.ToLower(operateType)]
	if has {
		return l
	}
	return LogOperateTypeNone
}

func GetLogOperateTitle(operateIndex LogOperateType) string {
	return logOperateTypeTitles[operateIndex]
}

type LogKindsItem struct {
	Title string `json:"title"`
	Name  string `json:"name"`
}

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
