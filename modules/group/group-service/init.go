package group_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	group := newCommonGroupService()

	bean.Injection(&group)
}
