package dynamic_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	ds := newDynamicService()
	bean.Injection(&ds)
}
