package user_model

import (
	"github.com/eolinker/apinto-dashboard/modules/user/user-entry"
)

type UserInfo struct {
	*user_entry.UserInfo
	OperateEnable bool
	Operator      string
}
