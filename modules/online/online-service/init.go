package online_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	resetOnline := newResetOnline()
	bean.Injection(&resetOnline)
}
