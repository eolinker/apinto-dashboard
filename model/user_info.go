package model

import "github.com/eolinker/apinto-dashboard/entry"

type UserInfo struct {
	*entry.UserInfo
	OperateEnable bool
	Operator      string
}
