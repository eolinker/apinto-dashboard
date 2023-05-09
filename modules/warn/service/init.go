package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	warnStrategy := newWarnStrategyService()
	warnHistory := newWarnHistoryService()
	bean.Injection(&warnStrategy)
	bean.Injection(&warnHistory)
}
