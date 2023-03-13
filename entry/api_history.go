package entry

import "time"

type ApiHistory struct {
	Id          int
	ApiId       int
	NamespaceId int
	OldValue    ApiHistoryInfo
	NewValue    ApiHistoryInfo
	OptType     OptType //1新增 2修改 3删除
	OptTime     time.Time
	Operator    int
}

type ApiHistoryInfo struct {
	Api    API              `json:"api"`
	Config APIVersionConfig `json:"config"`
}
