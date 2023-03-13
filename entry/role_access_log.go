package entry

import "time"

type AccessListLog struct {
	AccessIds []string `json:"access_ids,omitempty"` //权限数组
	Module    string   `json:"module,omitempty"`     //模块目录
}

type RoleAccessLog struct {
	Id       int
	Operator int //操作者ID
	RoleID   int //角色ID
	OldValue AccessListLog
	NewValue AccessListLog
	OptType  int
	OptTime  time.Time //创建时间
}
