package user_service

import "github.com/eolinker/eosc/common/bean"

func init() {

	userInfo := newUserInfoService()
	bean.Injection(&userInfo)
}
