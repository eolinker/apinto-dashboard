package navigation_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	service := newNavigationService()
	bean.Injection(&service)
}
