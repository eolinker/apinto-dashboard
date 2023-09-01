package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	infoService := newUserInfoService()
	bean.Injection(&infoService)
	userInfo := newUserInfoIdCache()
	userNameInfo := newUserInfoNameCache()
	session := newSessionCache()
	bean.Injection(&userInfo)
	bean.Injection(&userNameInfo)
	bean.Injection(&session)

}
