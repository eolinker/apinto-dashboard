package user_model

import (
	"github.com/eolinker/apinto-dashboard/entry/user-entry"
)

type UserInfo struct {
	*user_entry.UserInfo
	OperateEnable bool
	Operator      string
}
