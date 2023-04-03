package service

import "github.com/eolinker/eosc/common/bean"

func init() {
	infoService := newUserInfoService()
	bean.Injection(&infoService)
}
